Certainly! Here’s a comprehensive `README.md` file for both the Angular UI and the Go backend, including setup and configuration instructions:

---

# Shadowkeep Faucet

This project consists of two main components: the **Shadowkeep Faucet UI** (built with Angular 17) and the **Shadowkeep Faucet Backend** (built with Go). The UI allows users to claim test ETH on the Shadowkeep blockchain, while the backend handles token distribution, wallet interactions, and API management.

## Table of Contents

1. [Shadowkeep Faucet UI](#shadowkeep-faucet-ui)
2. [Shadowkeep Faucet Backend](#shadowkeep-faucet-backend)
3. [Setup and Configuration](#setup-and-configuration)
4. [Deployment](#deployment)
5. [Troubleshooting](#troubleshooting)
6. [License](#license)
7. [Contributing](#contributing)
8. [Acknowledgments](#acknowledgments)

## Shadowkeep Faucet UI

### Description

The Shadowkeep Faucet UI is a web application built with Angular 17 that allows users to connect their wallets and request test tokens from the Shadowkeep faucet. The interface is designed to be user-friendly and responsive.

### Features

- Wallet connection
- Token request functionality
- Display of network details
- Click-to-copy functionality for network information

### Build Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/KristiforMilchev/faucet-web
   cd shadowkeep-faucet-ui
   ```

2. **Install Dependencies**

   ```bash
   npm install
   ```

3. **Update Environment Configuration**

   Configure the environment settings in `src/environments/environment.ts` and `src/environments/environment.prod.ts` to match your backend API URL.

   ```typescript
   export const environment = {
     production: false,
     apiUrl: 'http://localhost:4445' // Replace with your backend URL
   };
   ```

4. **Run the Application**

   ```bash
   ng serve
   ```

5. **Build for Production**

   ```bash
   ng build --prod
   ```

## Shadowkeep Faucet Backend

### Description

The Shadowkeep Faucet Backend is a Go-based service that manages the distribution of test tokens on the Shadowkeep blockchain. It handles API requests from the UI, processes token claims, and manages wallet interactions.

### Features

- Token distribution
- API for wallet interactions
- JWT authentication
- PostgreSQL database integration

### Setup and Configuration

1. **Clone the Repository**

   ```bash
   git clone https://github.com/KristiforMilchev/faucet-backend
   cd shadowkeep-faucet-backend
   ```

2. **Build the Backend**

   Ensure you have Go installed. Build the Go backend using:

   ```bash
   go build -o faucet-backend
   ```

3. **Setup the PostgreSQL Database**

   Before starting the backend service, you need to create the database schema using the provided `db.sql` file.

   - **Install PostgreSQL**: If you haven't already, install PostgreSQL on your machine.
   - **Create the Database**:

     ```bash
     psql -U postgres -c "CREATE DATABASE faucet;"
     ```

   - **Apply the Schema**:

     ```bash
     psql -U postgres -d faucet -f db.sql
     ```

4. **Configure the Backend**

   Update the `settings.json` file with your specific configuration details. Ensure the file is properly set up with the connection string, JWT secret, and other configurations.

   ```json
   {
     "ConnectionString": "your_postgres_connection_string",
     "jwt-key": "your_jwt_key",
     "jwt-issuer": "your_jwt_issuer",
     "Port": ":4445",
     "Ledger": "your_ledger_private_key",
     "LedgerPublic": "your_ledger_public_address",
     "RPC": "your_rpc_url"
   }
   ```

5. **Start the Backend Service**

   Run the backend service using:

   ```bash
   ./faucet-backend
   ```

   Ensure that the backend is running and accessible at the configured port.

## Deployment

### UI Deployment

1. Build the Angular application using `ng build --prod`.
2. Deploy the contents of the `dist/` folder to your hosting service.

### Backend Deployment

1. Build the Go backend as described in the setup steps.
2. Deploy the `faucet-backend` binary to your server.
3. Ensure that PostgreSQL and the Go backend are properly configured and running.

## Troubleshooting

- **CORS Issues**: Ensure your Go backend is configured to allow CORS requests from your Angular app.
- **API Connection**: Verify that the API URL in your Angular environment configuration matches the backend's URL.
- **Database Issues**: Ensure PostgreSQL is running and the database schema is correctly applied.
- **JWT Issues**: Verify that the JWT secret in `settings.json` matches what the UI expects.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contributing

Contributions are welcome! Please fork the repository, make changes, and submit a pull request for review.

## Acknowledgments

- **Angular**: For the powerful front-end framework.
- **Go**: For the backend implementation.
- **PostgreSQL**: For the database management system.
- **MetaMask**: For Ethereum wallet integration.

---

Feel free to adjust any sections to better fit your specific project needs or additional details you might want to include.
