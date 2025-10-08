## The Problem with Sending Email in the Request Body

  Let's imagine we trust the email sent in the request body. The flow would be:
   1. A user logs into your dashboard with their Google account, say user@dlsu.edu.ph.
   2. The frontend gets a valid Google ID token for user@dlsu.edu.ph.
   3. The frontend sends a request to /request-key with:
       * Header: Authorization: Bearer <valid_token_for_user@dlsu.edu.ph>
       * Body: {"email": "user@dlsu.edu.ph", "project": "My Project"}

  This works fine for a normal user. But what if a malicious user, attacker@gmail.com, wants to generate a key for a legitimate R&D
  member, victim@dlsu.edu.ph?

  The attacker could do this:
   1. Log in to your dashboard with their own Google account, attacker@gmail.com.
   2. Get a valid Google ID token for attacker@gmail.com.
   3. Use a tool like curl to send a request to your API with:
       * Header: Authorization: Bearer <valid_token_for_attacker@gmail.com>
       * Body: {"email": "victim@dlsu.edu.ph", "project": "Evil Project"}

  Your server would see a valid Google token (for the attacker) and let the request through. The handler would then read the
  victim's email from the request body and issue a valid API key for them. The attacker now has an API key for a legitimate,
  privileged user.

  The Secure Way: Email from the ID Token (My Implementation)

  My proposed implementation prevents this attack entirely. Here's why it's better:

   1. The Google ID token is a cryptographically signed, tamper-proof "digital passport". The email claim inside it is the single
      source of truth for the user's identity.
   2. The middleware validates this passport. If it's authentic, it extracts the email directly from it.
   3. The handler only uses the email from the token. It completely ignores any email that might be sent in the body.

  In this secure model, the attacker can only ever generate an API key for their own email (attacker@gmail.com), because that's the
  email inside the token they possess. They can't generate a key for victim@dlsu.edu.ph because they don't have the victim's token.

  In short, by taking the email from the token, we guarantee that the user requesting the key is the actual owner of that email
  address.

  This is the standard, secure practice for this kind of authentication flow.
