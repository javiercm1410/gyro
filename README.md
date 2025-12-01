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
- [Roadmap](#roadmap)
- [Tasks](#tasks)
- [📄 License](#-license)

## 📘 Description

Gyro CLI is a command-line tool designed to list and rotate AWS access keys and users.

### Features

- **List:** List access keys and users.

## 🚢 Installation

To install gyro CLI from the source, follow these steps:

```bash
curl -s -L https://github.com/javiercm1410/gyro/releases/download/v0.3.0/gyro-darwin-x64.tar.gz | tar xz
chmod +x gyro
sudo mv gyro /usr/local/bin
```

## 🔧 Usage

### Commands

users: List AWS expired login Profiles
keys: List AWS expire keys

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

## Roadmap

- Rotation
- Delete only feature
- If 2user has two keys add confirmation
- User should provide temp pass
- Store/send results
- Add config file
- Write tests
- Notify with slack
- Mobile App

## Tasks

- Remove users without key from the output
- Check all commands
- Init rotation efforts

## 📄 License

This project is licensed under the MIT License. See the LICENSE file for details.


Manage this error if no login profile:
./gyro users -f json -u rotate



Manage error if user doesn't exist:
./gyro users -f json -u rotate

Check command and flags

Rotate

when user have no key:
./gyro keys -f json -u rotate
