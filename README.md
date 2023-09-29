# Lazy Twitter Clone
- Implemented with GO and Javascript
- Application made to showcase a couple of learnings from studying GO

# Points of improvement:
- Security
    - HTTPS configuration
    - Proper auth provider (oauth perhaps?), if not hashing of passwords
    - Move endpoints defined in JS source to a secure storage probably localStorage
- Database
    - ORM (using built in golang sql libraries for simplicity)
    - Decouple service layer from azure repo
- Libraries
    - Use of an HTTP framework
        - App uses the built in net/http library
            - This is for simplicity's sake and because the app is supposed to be for demo purposes
    - Use of a validator library for mapping between structs and request input
- UI Framework
    - Javascript framework, in the interest of time I opted to use vanilla JS
    - Preferably shouldve been built not as a SPA like application
    - A bundler for minifying our assets (most especially the JS)
- Tests
    - A test framework of some sort
- Architecture
    - As the app grows, move away from this MVC-esque approach to a microservice one where ideally the tweets is on its own
    - Usage of CDNs in the future for media (e.g. images, videos)
    - Update database to use a mix of RDBMS and NonRDBMS to cater rendering of tweet IDs on the server side when generating the feed, also ideally a pubsub model for updating said feed

# Features implemented
- Login
    - 1 day session token length
- Logout
- Display tweets of people you follow
    - With pagination
- Display your own tweets in your own feed
- Posting your own tweet
    - After posting you can view it in your feed
- Following users not in your follow list through the discover feed
    - 'Get another tweet' button for polling a suggested tweet from a user you don't follow

# Missing Features for the future:
- Viewing your own tweets via your profile 
    - Currently you can view your own tweets along with users you follow on the same feed page
    - The backend code for this already exists
- Deleting tweets
    - Backend code for this exists
- Following and unfollowing other users
    - Backend code for this exists
- Multi-page application
    - Currently implemented as a SPA for simplicity's sake
- Retweeting
- Pubsub for polling newly posted tweets as other users create more
- Favoriting tweets
- Uploading media