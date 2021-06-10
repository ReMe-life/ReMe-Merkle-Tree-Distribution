# ReMe-Merkle-Tree-Distribution

## ENV Preset

````
DB_HOST - Database host
DB_PORT - Database port
DB_USER - Database username
DB_PASS - Database password
DB_NAME - Database name
NODE_URL - Infura API Key
SECRET_KEY - Ethereum private key /this address should be the owner of the MerkleAirDrop contract/
CONTRACT_ADDRESS - AirDrop Contract address
SAVE_PERIOD - Saving period in seconds. How often the api will update merkle tree roothash.
API_PORT - Api port```

````

## Build & Deploy

docker build -t verifiable-merkle-tree:1.0 .
deploy on any preferred cloud service platform with sql database
`!IMPORTANT! The sql server should be configured to "Single zone"`
