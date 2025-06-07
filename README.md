# How to use kubb

## Follow these steps :
```sh
# Step 1: Clone the repository using the project's Git URL.
git clone [<git url>](https://github.com/Surasan01/courtopia-reserve.git)

# Step 2: Add .env file in backend/cmd/server.
MONGO_URI= <Your mongo url>
PORT= <Port>
JWT_SECRET= <your secret keyyy>
ENVIRONMENT=development

# Step 3: Run backend.
cd courtopia-reserve
cd backend
cd cmd
cd server
go run .

# Step 4: Run forntend.
cd courtopia-reserve
npm run dev

```
