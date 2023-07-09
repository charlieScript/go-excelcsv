# GBEDU

Easily convert Excel/CSV data to JSON format. 


## Features
- Convert Excel (.xlsx) files to JSON format
- Convert CSV (.csv) files to JSON format

## Installation
To install GBEDU, download the latest version for your operating system. It is available on Mac, Windows, and Linux.
Install the `.exe`, `.deb` and `.dmg` files for Windows, Linux and Mac OSes respectively.

## How To Use
- Select the source file format (Excel/CSV).
- Check the `Header` button if the first row of the data should be set as header.
- Select the source file from the file dialog.
- Copy or download the JSON data.

## Build/Run
To build or run `GBEDU` locally on your machine, first install [Golanf](https://go.dev/).  

Then:
1. Clone the repository  

    ```shell
    git clone https://github.com/charlieScript/go-excelcsv
    ```
2. Install dependencies  

    ```shell
    go install .
    ```
3. Run the app  

    ```shell
    go run .
    ```  

   or build the app  

    ```shell
    go build .
    ```

inspired by: https://github.com/Johnkayode/kodiak