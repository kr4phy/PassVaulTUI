# PassVaulTUI
PassVaulTUI is a lightweight, terminal-based password manager written in Go.
Manage your passwords securely anywhere with a single binary under 10MB, without complex configuration.

## Key Features
- Credential Management: Effortlessly create (Save), modify (Edit), and remove (Delete) password entries.
- Secure Password Generator: Generate strong, random passwords instantly.
- Ultra-lightweight: Distributed as a single binary under 10MB for fast and portable use.
- TUI Support: Features an intuitive interface optimized for terminal environments.
- Fully offline and local storage: No cloud dependencies, ensuring your data remains private and secure.
- AES-GCM Encryption: All data is encrypted using AES-GCM for robust security.

## Getting Started

### Prerequisites

To run or build this project, you need to have Go 1.25+ installed on your system.
Encrypted data is saved to `./data.bin`; if it does not exist, it is created automatically.

### Installation & Execution

#### Clone the Repository
```
git clone https://github.com/kr4phy/PassVaulTUI.git
cd PassVaulTUI
```

#### Run Directly
```
go run .
```

#### Build Binary
```
go build .
```

Then run `PassVaulTUI.exe`(Windows) or `PassVaulTUI`(Mac OS/Linux).

On first run, just type the password you want.
Then, `./data.bin` will be automatically generated.
If `./data.bin` already exists, enter the password you set.

All data is stored locally in `./data.bin`; deleting it permanently removes all stored data.
If you want to back up or restore data, copy `data.bin` from the binary's directory.
When repairing data, `data.bin` must be placed in the binary's directory.

CLI/Headless mode is not supported to avoid exposing the master password.

## License
This project is licensed under the BSD 3-Clause License. For more details, please see the LICENSE file.

## Contributing
Bug reports, feature suggestions, and Pull Requests (PR) are always welcome! If you'd like to help improve the project, please feel free to open an issue.