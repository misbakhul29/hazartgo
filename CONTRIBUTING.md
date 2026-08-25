# 🤝 Contributing to HazartGo

Thank you for your interest in contributing to **HazartGo**! We welcome contributions from developers of all skill levels.

---

## 🛠️ How to Contribute

### 1. Fork the Repository
Click the **Fork** button at the top right of this repository on GitHub to create your own copy.

### 2. Clone Your Fork
```bash
git clone https://github.com/YOUR_USERNAME/hazartgo.git
cd hazartgo
```

### 3. Create a Feature Branch
```bash
git checkout -b feature/my-cool-feature
```

### 4. Make Changes & Run Tests
Ensure all unit tests pass before committing:
```bash
go test ./...
```

### 5. Commit Your Changes
Use conventional commit style (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`):
```bash
git commit -m "feat: add new feature"
```

### 6. Push & Create a Pull Request (PR)
```bash
git push origin feature/my-cool-feature
```
Open a Pull Request on GitHub targeting the `dev` branch.

---

## 📜 Pull Request Guidelines
- All PRs must pass automated GitHub Actions CI builds & unit tests.
- Maintain standard Go code formatting (`gofmt`).
- Include unit tests for any new features or bug fixes.
