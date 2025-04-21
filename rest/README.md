# 📦 Package Watcher API V1

Welcome to the **Watcher API** documentation.

---

## Features 💎

- Packages
  - List
- Single package
  - Get data
  - Install
  - Delete (not done yet)
  - Initialize
  - Run
  - Shutdown

---

## 🌐 Base URL

```
http://localhost/api/v1
```

---

## 🚀 Endpoints

### 🔹 Get All Packages

**GET** `/package`

Retrieve the list of all available packages.

**Request Body**

```json
{}
```

---

### 🔹 Get Single Package

**GET** `/package/{name}`

Get information about a specific package.

**Path Parameter**

- `name` (string): Package name.

**Request Body**

```json
{}
```

---

### 🔹 Install Package

**POST** `/package/{name}`

Install a package.

**Path Parameter**

- `name` (string): Package name.

**Success Response**

```json
{
  "message": "Package successfully installed."
}
```

---

### 🔹 Delete Package

**DELETE** `/package/{name}`

Delete a package from the system.

**Path Parameter**

- `name` (string): Package name.

**Response:**

- No content (`204 No Content` if successful).

---

### 🔹 Initialize Package

**PATCH** `/package/{name}/init`

Initialize the package process.

**Path Parameter**

- `name` (string): Package name.

**Success Response**

```json
{
  "message": "Package initialized successfully."
}
```

---

### 🔹 Run Package

**PATCH** `/package/{name}/run`

Run the package.

**Path Parameter**

- `name` (string): Package name.

**Success Response**

```json
{
  "message": "Package is now running."
}
```

---

### 🔹 Shutdown Package

**PATCH** `/package/{name}/stop`

Shut down the package process.

**Path Parameter**

- `name` (string): Package name.

**Success Response**

```json
{
  "message": "Package process stopped successfully."
}
```

---

## 🛠️ Quick Usage Example

```bash
curl -X POST http://localhost/api/v1/package/my-package
```

---

## 🧠 Notes

- Package names **must not contain spaces**. Use hyphens (`-`) instead.
- All routes follow RESTful conventions.
- Error and success messages are returned as terminal message.
