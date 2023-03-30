// // Import the functions you need from the SDKs you need
// import { initializeApp } from "firebase/app";
// import { getAnalytics } from "firebase/analytics";
// import { getDatabase } from "firebase/database";
// // TODO: Add SDKs for Firebase products that you want to use
// // https://firebase.google.com/docs/web/setup#available-libraries

// // Your web app's Firebase configuration
// // For Firebase JS SDK v7.20.0 and later, measurementId is optional
// const firebaseConfig = {
//   apiKey: "AIzaSyD7wKLcgoR2QCNsKjPh14coLtjLg7g1fvM",
//   authDomain: "registrytest-acbf4.firebaseapp.com",
//   databaseURL: "https://registrytest-acbf4-default-rtdb.firebaseio.com",
//   projectId: "registrytest-acbf4",
//   storageBucket: "registrytest-acbf4.appspot.com",
//   messagingSenderId: "484449986816",
//   appId: "1:484449986816:web:2e72fb47d5ebfc990eabe5",
//   measurementId: "G-7KV1XBG7Q8"
// };

// // Initialize Firebase
// const app = initializeApp(firebaseConfig);
// // const analytics = getAnalytics(app);
// const database = getDatabase(app);

// export default database;

import firebase from 'firebase/compat/app';
import { initializeApp } from "firebase/app";
import { getDatabase, ref, set} from "firebase/database";

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
// const firebaseConfig = {
//   apiKey: "AIzaSyD7wKLcgoR2QCNsKjPh14coLtjLg7g1fvM",
//   authDomain: "registrytest-acbf4.firebaseapp.com",
//   databaseURL: "https://registrytest-acbf4-default-rtdb.firebaseio.com",
//   projectId: "registrytest-acbf4",
//   storageBucket: "registrytest-acbf4.appspot.com",
//   messagingSenderId: "484449986816",
//   appId: "1:484449986816:web:2e72fb47d5ebfc990eabe5",
//   measurementId: "G-7KV1XBG7Q8"
// };

const firebaseConfig = JSON.parse(process.env.FIREBASE_CONFIG);

const app = initializeApp(firebaseConfig);

const db = getDatabase(app);
  
// Push some test data to the "packages" table
set(ref(db, 'package/' + "react"), {
  "name": 'react',
  url: 'https://github.com/facebook/react',
  net_score: 8.3,
  ramp_up_score: 7.2,
  correctness_score: 9.1,
  bus_factor_score: 6.5,
  responsiveness_score: 8.9,
  license_score: 10
});

set(ref(db, 'package/' + "lodash"), {
  name: 'lodash',
    url: 'https://github.com/lodash/lodash',
    net_score: 7.2,
    ramp_up_score: 8.1,
    correctness_score: 8.3,
    bus_factor_score: 7.5,
    responsiveness_score: 9.2,
    license_score: 8.5
});

// const packagesRef = db.ref('packages');
// packagesRef.push({
//   name: 'react',
//   url: 'https://github.com/facebook/react',
//   net_score: 8.3,
//   ramp_up_score: 7.2,
//   correctness_score: 9.1,
//   bus_factor_score: 6.5,
//   responsiveness_score: 8.9,
//   license_score: 10
// });

// // Push another test package
// packagesRef.push({
//   name: 'lodash',
//   url: 'https://github.com/lodash/lodash',
//   net_score: 7.2,
//   ramp_up_score: 8.1,
//   correctness_score: 8.3,
//   bus_factor_score: 7.5,
//   responsiveness_score: 9.2,
//   license_score: 8.5
// });

// db.goOffline();
// process.exit(0);

export default db;