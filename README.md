# la-accessor-geocode

### how to call the api

Step 1. 

Curl an address like this
```
https://assessor.lacontroller.io/getcoords?address=200%20N%20MAIN%20ST
```

That's it!

### How to set this up yourself!

Step 1. Make a postgres database.

Step 2. Make the table using the script in maketable.sql

Step 3. COPY your CSV into the `parcels` table.

Step 4. Compile and deploy the API, in this case, we are using a VM and deploying it on port 8080 but feel free to use Google Kubernetes Engine or whatever.

