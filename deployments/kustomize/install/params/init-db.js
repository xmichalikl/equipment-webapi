const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.AMBULANCE_API_MONGODB_DATABASE
const collection = process.env.AMBULANCE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

function generateUUID() { // Public Domain/MIT
  var d = new Date().getTime();//Timestamp
  var d2 = (typeof performance !== 'undefined' && performance.now && (performance.now()*1000)) || 0;//Time in microseconds since page-load or 0 if unsupported
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      var r = Math.random() * 16;//random number between 0 and 16
      if(d > 0){ //Use timestamp until depleted
          r = (d + r)%16 | 0;
          d = Math.floor(d/16);
      } else { //Use microseconds since page-load if supported
          r = (d2 + r)%16 | 0;
          d2 = Math.floor(d2/16);
      }
      return (c==='x' ? r : (r&0x3|0x8)).toString(16);
  });
}

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
       print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

//insert sample data
let result = db[collection].insertMany([
    {
        "id": generateUUID(),
        "name": "Chirurgia",
        "equipment": [
          {
            "id": generateUUID(),
            "name": "Defibrilátor",
            "availability": "available",
            "lastInspectionDate": new Date("2024-05-22T08:17:23.950Z"),
            "technicalCondition": 3,
            "inspectionInterval": 6
          },
        ]
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);