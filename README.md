<!-- omit in toc -->
# Gyro CLI

<!-- omit in toc -->
## Contents

- [📘 Description](#-description)
  - [Features](#features)
- [🚢 Installation](#-installation)
- [🔧 Usage](#-usage)
  - [Commands](#commands)
  - [Examples](#examples)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## 📘 Description

Gyro CLI is a command-line tool designed to generate .env files using various providers, including Werf and 1Password.

### Features

- **Create .env with Werf:** Generate gyroronment files based on configurations managed by Werf.
- **Create .env with 1Password:** Securely generate gyroronment files using secrets stored in 1Password.

## 🚢 Installation

To install gyro CLI from the source, follow these steps:

```bash
curl -s -L https://github.com/javiercm1410/gyro/releases/download/v0.3.0/gyro-darwin-x64.tar.gz | tar xz
chmod +x gyro
sudo mv gyro /usr/local/bin
```

## 🔧 Usage

Once installed, you can use the gyro command to generate .env files based on Werf or 1Password.

### Commands

werf: Generate a .env file using Werf configurations.
1pass: Generate a .env file using 1Password secrets.
help: Display help information about any command.

### Examples

To list AWS access keys

```bash
./gyro keys
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (git checkout -b feature-branch).
3. Make your changes.
4. Commit your changes (git commit -m 'Add new feature').
5. Push to the branch (git push origin feature-branch).
6. Open a Pull Request.

## 📄 License

This project is licensed under the MIT License. See the LICENSE file for details.


ROtate Access KEY
Rotate PAssowrd
Password manager as provider
Use config file


Multiple Output types
Check Flags
Check write to file


they should be two command one for key
the other for pass

key command
user command

rotate will have flags


Old in red, close to old in yellow


Exit on failure


safe key on aws configure

use config file

add notify feature
