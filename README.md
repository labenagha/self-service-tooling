```markdown
# Self-Service Tooling App

This is a full-stack application for managing Terraform infrastructure runs and user authentication. It allows users to trigger Terraform actions like `Plan`, `Apply`, and `Destroy` through a responsive interface. The app is built with **Node.js**, **Express**, and **SQLite** for the backend, and vanilla **HTML**, **CSS**, and **JavaScript** for the ui.

---

## Features

- **User Authentication**: Register, login, and manage user sessions using JSON Web Tokens (JWT).
- **Trigger Terraform Runs**: Trigger `Plan`, `Apply`, and `Destroy` operations.
- **Run Logs**: View detailed logs of your Terraform runs.
- **Responsive Design**: Mobile-friendly interface using responsive CSS.
- **Easy Hosting**: Compatible with free hosting platforms like Heroku, Railway, and Render.

---

## Prerequisites

Before using this app, ensure you have the following installed:

- **Node.js** (v14+)
- **npm** (v6+)
- **Git**
- **SQLite**

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/self-service-tooling.git
cd self-service-tooling
```

### 2. Backend Setup

1. Navigate to the `backend` directory:

   ```bash
   cd backend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Create a `.env` file in the `backend` directory and add the following variables:

   ```env
   PORT=5000
   JWT_SECRET=your_jwt_secret
   ```

4. Start the backend server:

   ```bash
   npm run dev
   ```

   The server should now be running at `http://localhost:5000`.

---

### 3. ui Setup

1. Navigate to the `ui` directory:

   ```bash
   cd ../ui
   ```

2. Open the `ui/index.html` file in your browser.

---

## Directory Structure

```plaintext
self-service-tooling/
├── backend/
│   ├── controllers/      # Contains backend logic for authentication
│   ├── database/         # SQLite database and configuration
│   ├── models/           # User model for database queries
│   ├── routes/           # API routes for authentication and actions
│   ├── index.js          # Entry point for the backend server
│   ├── package.json      # Backend dependencies
├── ui/
│   ├── assets/           # Images, logos, and icons
│   ├── scripts/          # ui JavaScript (login.js, register.js, etc.)
│   ├── styles/           # ui CSS
│   ├── index.html        # Main app interface
│   ├── login.html        # Login page
│   ├── register.html     # Register page
│   ├── plan.html         # Terraform Plan page
│   ├── apply.html        # Terraform Apply page
```

---

## Usage

### User Authentication

1. Open the app in your browser.
2. Click **"Create Account"** to register a new user.
3. Login using your credentials to access the app.

### Trigger a Run

1. Navigate to the **"Trigger A Run"** section.
2. Select the desired Terraform action:
   - `Plan`: View the execution plan for your Terraform configuration.
   - `Apply`: Apply the Terraform configuration to provision resources.
   - `Destroy`: Tear down infrastructure managed by Terraform.
3. Follow the logs to track the status of the run.

---

## API Endpoints

### User Authentication

| Method | Endpoint       | Description                |
|--------|----------------|----------------------------|
| POST   | `/api/register` | Register a new user        |
| POST   | `/api/login`    | Login and receive a JWT    |

### Protected Routes

| Method | Endpoint          | Description                |
|--------|-------------------|----------------------------|
| GET    | `/api/protected`  | Access protected content   |

---

## Responsive Design

The app is fully responsive, ensuring it works across desktop, tablet, and mobile devices. Key breakpoints include:

- **768px**: Adjusts layout for tablets.
- **480px**: Optimizes layout for smaller screens.

---

## Deployment

### Deploying the Backend

1. Use a free hosting platform like [Render](https://render.com) or [Railway](https://railway.app).
2. Push the `backend` directory to your hosting provider.
3. Set environment variables (`PORT`, `JWT_SECRET`) in the hosting provider's dashboard.

### Deploying the ui

1. Use a static hosting service like [Netlify](https://www.netlify.com) or [Vercel](https://vercel.com).
2. Upload the `ui` directory.
3. Configure the backend API URL if necessary.

---

## Built With

- **Backend**: Node.js, Express, SQLite
- **ui**: HTML, CSS, JavaScript
- **Authentication**: JSON Web Tokens (JWT)

---

## Example `.env` File

For your backend server, create a `.env` file with the following content:

```env
PORT=5000
JWT_SECRET=your_jwt_secret
```

---

## Contributing

1. Fork the repository.
2. Create a new feature branch:

   ```bash
   git checkout -b feature-name
   ```

3. Commit your changes:

   ```bash
   git commit -m "Add new feature"
   ```

4. Push to the branch:

   ```bash
   git push origin feature-name
   ```

5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contact

For issues or feature requests, please open an issue on [GitHub](https://github.com/labenagha/self-service-tooling).
