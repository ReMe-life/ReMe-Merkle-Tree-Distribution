package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LimeChain/merkletree"
	"github.com/LimeChain/merkletree/memory"
	"github.com/LimeChain/merkletree/postgres"
	merkleRestAPI "github.com/LimeChain/merkletree/restapi/baseapi"
	validateAPI "github.com/LimeChain/merkletree/restapi/validateapi"
	env "github.com/Netflix/go-env"
	"github.com/ReMe-life/ReMe-Merkle-Tree-Distribution/saver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Environment struct {
	Blockchain struct {
		NodeURL         *string `env:"NODE_URL"`
		Secret          *string `env:"SECRET_KEY"`
		ContractAddress *string `env:"CONTRACT_ADDRESS"`
		SavePeriod      int     `env:"SAVE_PERIOD"`
	}

	Database struct {
		Host    *string `env:"DB_HOST"`
		Port    int     `env:"DB_PORT"`
		User    *string `env:"DB_USER"`
		Pass    *string `env:"DB_PASS"`
		DBName  *string `env:"DB_NAME"`
		DBExtra *string `env:"DB_EXTRA"`
	}

	Port int `env:"API_PORT"`

	Extras env.EnvSet
}

var treeLen int

type savedResponse struct {
	merkleRestAPI.MerkleAPIResponse
	Length int `json:"lastSavedIndex"`
}

func getSaved(tree merkletree.ExternalMerkleTree) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tree.Length() == 0 {
			render.JSON(w, r, savedResponse{merkleRestAPI.MerkleAPIResponse{true, ""}, -1})
			return
		}
		render.JSON(w, r, savedResponse{merkleRestAPI.MerkleAPIResponse{true, ""}, treeLen})
		return
	}
}

func createAndStartAPI(tree merkletree.ExternalMerkleTree, port int) {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Compress(6, "gzip"),
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.NoCache,
	)

	router.Route("/v1", func(r chi.Router) {
		treeRouter := chi.NewRouter()
		treeRouter = merkleRestAPI.MerkleTreeStatus(treeRouter, tree)
		treeRouter = merkleRestAPI.MerkleTreeInsert(treeRouter, tree)
		treeRouter = merkleRestAPI.MerkleTreeHashes(treeRouter, tree)
		// treeRouter = merkleRestAPI.MerkleTreeRawInsert(treeRouter, tree)
		treeRouter = validateAPI.MerkleTreeValidate(treeRouter, tree)

		treeRouter.Get("/saved", getSaved(tree))
		r.Mount("/api/merkletree", treeRouter)
	})

	fmt.Printf("Starting REST Api at port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func createSaver(tree merkletree.FullMerkleTree, nodeUrl, privateKeyHex, contractAddress string, periodSeconds int) {
	treeSaver, err := saver.NewSaver(
		nodeUrl,
		privateKeyHex,
		contractAddress,
		tree)
	if err != nil {
		panic(err)
	}

	go func() {
		len := 0
		timeout := time.Duration(periodSeconds) * time.Second
		for {
			tree.Recalculate()
			savedRoot, err := treeSaver.FetchRoot()
			if err != nil {
				fmt.Println("ERR: Could not fetch root")
				fmt.Println(err.Error())
				time.Sleep(timeout)
				continue
			}

			if savedRoot == tree.Root() {
				fmt.Printf("Same root (%v) found in the chain. Skipping this iteration\n", savedRoot)
				treeLen = tree.Length()
				time.Sleep(timeout)
				continue
			}

			treeLen = tree.Length()
			if treeLen > len {
				fmt.Printf("Submitting new tree root to the chain (%v)\n", tree.Root())
				tx, err := treeSaver.TriggerSave()
				if err != nil {
					fmt.Println("ERR: Could not save the tree root")
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Transaction (%v) mined\n", tx.TxHash.Hex())
					fmt.Printf("Gas used (%v)\n", tx.GasUsed)
					len = treeLen
				}
			} else {
				fmt.Println("No changes to submit. Skipping this iteration")
			}
			time.Sleep(timeout)
		}
	}()

	fmt.Printf("Started saver on %v seconds\n", periodSeconds)
	fmt.Printf("Node url %v\n", nodeUrl)
	fmt.Printf("Verifier contract address %v \n", contractAddress)
}

func loadPostgreTree(connStr string) merkletree.FullMerkleTree {
	tree := postgres.LoadMerkleTree(memory.NewMerkleTree(), connStr)
	fmt.Printf("Merkle tree loaded. Length : %v, Root : %v\n", tree.Length(), tree.Root())
	return tree
}

func main() {

	var connStr, nodeURL, privateKeyHex, contractAddress string
	var period, port int

	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}
	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", *environment.Database.Host, environment.Database.Port, *environment.Database.User, *environment.Database.Pass, *environment.Database.DBName)
	nodeURL = *environment.Blockchain.NodeURL
	privateKeyHex = *environment.Blockchain.Secret
	contractAddress = *environment.Blockchain.ContractAddress
	period = environment.Blockchain.SavePeriod
	port = environment.Port

	tree := loadPostgreTree(connStr)

	createSaver(tree, nodeURL, privateKeyHex, contractAddress, period)

	createAndStartAPI(tree, port)
}
