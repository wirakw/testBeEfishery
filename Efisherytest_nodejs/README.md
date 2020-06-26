# eFishery Test Rest Api Auth With Node
# singgle running program
   1. create file .env berdasarkan .env_example
   2. ganti data environtmennya khusus untuk database_url nya saja mengikuti database_url yang tersedia berserta kaidah penulisan database_urlnya
   
   3. run 
      ```shell
         $ npm install
      ```
   4.  ```shell
         $ npm run create-tables
       ```
       untuk migration tabel users jika belom ada

   5.  run
       ```shell
         $ npm run start
      ```
# running all with docker composer
   1. run 
      ```shell
         $ docker-compose up -d
      ```