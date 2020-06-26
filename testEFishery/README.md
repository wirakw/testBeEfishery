# testEFishery
# singgle running program
   1. set environment variable menggunakan expoxe mengikuti .env_example (library get file godoc yang saya gunakan tidak suport untuk mengambil data environment dari file .env dan ENVIRONMENT variable dari docker composer sekaligus, jadi saya memutuskan untuk memilih running di docker composer)
   2. ganti data environtmennya sesuai postgres yang tersedia
   3. run 
      ```shell
         $ go run .
      ```
catatan : singgle running lebih menyebabkan rentan error
# running all with docker composer
   1. run 
      ```shell
         $ docker-compose up -d
      ```