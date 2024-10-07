const express = require("express");
const passport = require("passport");
const GoogleStrategy = require("passport-google-oidc");
const db = require("../db");

const router = express.Router();

/* NOTES:
 
profile: {
    id: '<unique>',
    displayName: 'Edwin Sadiarin Jr.',
    name: { 
        familyName: 'Sadiarin Jr.',
        givenName: 'Edwin' 
    },
    emails: [ { value: 'edwin_sadiarinjr@dlsu.edu.ph' } ]
}


 * */

passport.use(
    new GoogleStrategy(
        {
            clientID: process.env["GOOGLE_CLIENT_ID"],
            clientSecret: process.env["GOOGLE_CLIENT_SECRET"],
            callbackURL: "/oauth2/redirect/google",
            scope: ["email", "profile"],
        },
        // callback
        function verify(issuer, profile, cb) {
            // log the profile data for debugging purposes
            console.log("Authenticated user profile:", profile);
            console.log("User's display name:", profile.displayName);
            console.log("User's email:", profile.emails ? profile.emails[0].value : "No email provided");

            // store the user in the database after successful login
            db.get(
                "SELECT * FROM federated_credentials WHERE provider = ? AND subject = ?",
                [issuer, profile.id],
                function (err, row) {
                    if (err) {
                        return cb(err);
                    }
                    if (!row) {
                        db.run("INSERT INTO users (name) VALUES (?)", [profile.displayName], function (err) {
                            if (err) {
                                return cb(err);
                            }

                            const id = this.lastID;
                            db.run(
                                "INSERT INTO federated_credentials (user_id, provider, subject) VALUES (?, ?, ?)",
                                [id, issuer, profile.id],
                                function (err) {
                                    if (err) {
                                        return cb(err);
                                    }
                                    const user = {
                                        id: id,
                                        name: profile.displayName,
                                    };
                                    return cb(null, user);
                                },
                            );
                        });
                    } else {
                        db.get("SELECT * FROM users WHERE id = ?", [row.user_id], function (err, row) {
                            if (err) {
                                return cb(err);
                            }
                            if (!row) {
                                return cb(null, false);
                            }
                            return cb(null, row);
                        });
                    }
                },
            );
        },
    ),
);

passport.serializeUser(function (user, cb) {
    process.nextTick(function () {
        cb(null, { id: user.id, username: user.username, name: user.name });
    });
});

passport.deserializeUser(function (user, cb) {
    process.nextTick(function () {
        return cb(null, user);
    });
});

router.post("/logout", function (req, res, next) {
    req.logout(function (err) {
        if (err) {
            return next(err);
        }
        // maybe also destory session?
        res.redirect("/");
    });
});

// renders page (web)
router.get("/login", function (req, res, next) {
    res.render("login");
    // passport.authenticate("google");
});

router.get("/login/google", passport.authenticate("google"));

router.get(
    "/oauth2/redirect/google",
    passport.authenticate("google", {
        // successRedirect: "/",
        successRedirect: "/web", // for testing with web frontend
        failureRedirect: "/login",
    }),
);

module.exports = router;
