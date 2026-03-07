# 🤝 Contributing to qwen-pet

We welcome contributions to qwen-pet! To make sure the process goes smoothly, please follow these guidelines:

## 📋 Code of Conduct

Please note that all participants in our project are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md). Make sure to review it before contributing.

## 🛠 How to Contribute

1. **Fork the repository**:
   Fork the project to your GitHub account using the GitHub interface.

2. **Create a new branch**:
   Use a descriptive branch name for your feature or bugfix:

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**:
   Implement your feature or fix the bug in your branch. Make sure to include tests where applicable and follow coding standards.

4. **Test your changes**:

   ```bash
   go fmt ./...
   go vet ./...
   go test -v ./...
   ```

5. **Commit your changes**:
   Use conventional commit messages:

   ```bash
   git commit -m "feat(scope): description of changes"
   ```

6. **Push your changes**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Submit a Pull Request**:
   Create a pull request into the `dev` branch, detailing the changes you've made. Link any issues your changes resolve and provide context.

## 📑 Guidelines for Contributions

- **Format your code** with `gofmt` before submitting a pull request.
- **Vet your code** with `go vet ./...` to catch common issues.
- Ensure **test coverage** for your code. Uncovered code may delay the approval process.
- Write clear, concise **commit messages** following [conventional commits](https://www.conventionalcommits.org/).

Thank you for helping improve!

---

## 📜 License

2026 - This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). You are free to use, modify, and distribute this software under the terms of the GPL-3.0 license. For more details, please refer to the [LICENSE](LICENSE) file included in this repository.
