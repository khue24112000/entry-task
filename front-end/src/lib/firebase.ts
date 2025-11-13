import {
  initializeApp,
  getApps,
  getApp,
} from "firebase/app";
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "firebase/storage";

const firebaseConfig = {
  apiKey: "AIzaSyAI23M9oDqpwbDkSQS1lnd09tvsCQHJrZA",
  authDomain: "shop-ae551.firebaseapp.com",
  projectId: "shop-ae551",
  storageBucket: "shop-ae551.appspot.com",
  messagingSenderId: "160755370410",
  appId: "1:160755370410:web:da81c2ef17dd22c6fd03d3",
  measurementId: "G-0FGJJBSH46",
};

export const uploadImage = async (file: File) => {
  const app = !getApps().length
    ? initializeApp(firebaseConfig)
    : getApp();

  const storage = getStorage(app);

  const uniqueFilename = Date.now() + "_" + file.name;

  const storageRef = ref(
    storage,
    "images/" + uniqueFilename
  );

  try {
    const snapshot = await uploadBytes(storageRef, file);
    const downloadURL = await getDownloadURL(snapshot.ref);
    return downloadURL;
  } catch (error) {
    console.error("Catch some error", error);
  }
};
