import React, {useState} from "react";
import { signInWithGoogle } from "../firebase";

const SignIn = () => {
    const [error] = useState(null);

  return (
    <div className="mt-8">
      <h1 className="text-3xl mb-2 text-center font-bold">Top</h1>
      <div className="border border-blue-400 mx-auto w-11/12 md:w-2/4 rounded py-8 px-4 md:px-8">
        {error !== null && <div className = "py-4 bg-red-600 w-full text-white text-center mb-3">{error}</div>}
        <button
          className="bg-red-500 hover:bg-red-600 w-full py-2 text-white"
          onClick={() => {
            signInWithGoogle();
          }}
        >
          Sign in with Google
        </button>
      </div>
    </div>
  );
};

export default SignIn;