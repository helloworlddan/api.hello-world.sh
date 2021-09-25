import firebase from "firebase/app";
import "firebase/auth";

const firebaseConfig = {
    apiKey: "",
    authDomain: "hwsh-api.firebaseapp.com",
    databaseURL: "https://hwsh-api.firebaseio.com",
    projectId: "hwsh-api",
    storageBucket: "hwsh-api.appspot.com",
    messagingSenderId: "",
    appId: ""
};
firebase.initializeApp(firebaseConfig);

export const auth = firebase.auth();

const googleProvider = new firebase.auth.GoogleAuthProvider();
export const signInWithGoogle = () => {
  auth.signInWithPopup(googleProvider);
};

const githubProvider = new firebase.auth.GithubAuthProvider();
export const signInWithGithub = () => {
  auth.signInWithPopup(githubProvider);
};
