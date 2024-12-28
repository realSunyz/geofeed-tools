![geofeed-validator](https://socialify.git.ci/realSunyz/geofeed-validator/image?description=1&descriptionEditable=&font=Jost&language=1&name=1&owner=1&pattern=Circuit%20Board&theme=Auto)

This utility is used to validate the format of geofeed files.

> [!NOTE]
> The project is still under development. The code on the main branch is stable and usable.

## Usage

### Compile

```bash
# Clone the repository
git clone https://github.com/realSunyz/geofeed-validator.git

# Navigate to the directory
cd geofeed-validator/src

# Synchronize and clean up dependencies
go mod tidy

# Build the executable file
go build -o geofeed-validator
```

### Run
```bash
./geofeed-validator <path_to_geofeed.csv>
```

## Contributing

Issues and Pull Requests are definitely welcome!

Please make sure you have tested your code locally before submitting a PR.

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/realSunyz/geofeed-validator/blob/main/LICENSE) file for details.