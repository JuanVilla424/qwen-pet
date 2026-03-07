## 🔨 Installation

1. **Clone the Repository**

   ```bash
   git clone --recurse-submodules https://github.com/JuanVilla424/qwen-pet.git
   ```

2. **Navigate to the Project Directory**

   ```bash
   cd qwen-pet
   ```

3. **Initialize Submodules** (if cloned without `--recurse-submodules`)

   ```bash
   git submodule update --init --recursive
   ```

4. **Set Up Python Virtual Environment** (CICD tooling)

   ```bash
   python -m venv venv
   source venv/bin/activate
   pip install --upgrade pip
   pip install poetry
   poetry lock
   poetry install
   ```

5. **Install Pre-Commit Hooks**

   ```bash
   source venv/bin/activate
   pre-commit install
   pre-commit install -t pre-commit
   pre-commit install -t pre-push
   ```

6. **Configure Environment**

   ```bash
   cp .env.template .env
   # Edit .env with your credentials
   ```

7. **Build**

   ```bash
   go build ./cmd/...
   ```
