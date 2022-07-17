import React, { useState } from "react";
import { Link } from "@reach/router";
import { signInWithGoogle } from "../firebase";

const SignUp = () => {
  const [error] = useState(null);

  return (
    <div className="mt-8">
      <h1 className="text-3xl mb-2 text-center font-bold">Sign Up</h1>
      <div className="border border-blue-400 mx-auto w-11/12 md:w-2/4 rounded py-8 px-4 md:px-8">
        {error !== null && (
          <div className="py-4 bg-red-600 w-full text-white text-center mb-3">
            {error}
          </div>
        )}
        <p className="text-center my-3">or</p>
        <button
          onClick={() => {
            try {
              signInWithGoogle();
            } catch (error) {
              console.error("Error signing in with Google", error);
            }
          }}
          className="bg-red-500 hover:bg-red-600 w-full py-2 text-white"
        >
          Sign In with Google
        </button>
        <p className="text-center my-3">
          Already have an account?{" "}
          <Link to="/" className="text-blue-500 hover:text-blue-600">
            Sign in here
          </Link>{" "}
        </p>
      </div>
    </div>
  );
};

export default SignUp;
