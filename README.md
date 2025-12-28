# PassVaulTUI
PassVaulTUI is a lightweight, terminal-based password manager written in Go.
Manage your passwords securely anywhere with a single binary under 10MB, without the need for complex configurations.

## Key Features
- Credential Management: Effortlessly create (Save), modify (Edit), and remove (Delete) password entries.
- Security Generator: Generate strong, random passwords instantly.
- Ultra-lightweight: Distributed as a single binary under 10MB for fast and portable use.
- TUI Support: Features an intuitive interface optimized for terminal environments.
- Fully offline and local storage: No cloud dependencies, ensuring your data remains private and secure.
- AES-GCM Encryption: All data is encrypted using AES-GCM for robust security.

## Getting Started
Prerequisites
To run or build this project, you need to have Go installed on your system.

## Installation & Execution
Clone the Repository
```
git clone https://github.com/kr4phy/PassVaulTUI.git
cd PassVaulTUI
```
Run Directly
```
go run .
```
Build Binary
```
go build -o passvaultui .
./passvaultui
```

## License
This project is licensed under the BSD 3-Clause License. For more details, please see the LICENSE file.

## Contributing
Bug reports, feature suggestions, and Pull Requests (PR) are always welcome! If you'd like to help improve the project, please feel free to open an issue.