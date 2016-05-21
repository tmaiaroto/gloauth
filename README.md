# GLOAuth

Go Lambda OAuth ("glow auth") is a serverless authentication system relying on AWS (Lambda, DynamoDB, Cognito, SES, API Gateway).

It's goal is to build a simple, yet flexible, zero maintenance authentication system for all of your application needs. 
Why another one? Well, there's a lot of good auth packages for Node.js (Passport for example) but I haven't find any for
Golang nor do I find many that leverage AWS Lambda (there is LambdAuth, again in Node.js). At least not at the time I needed
this. If I did, I wouldn't be building this, trust me.

### Why Go?

Well, LambdAuth wasn't exactly turn-key. The setup process wasn't working out for me and I quickly realized that not only
was it not going to be as "easy" as I thought, but it was also done with Node.js so it's slower. Not a lot of course, but
when we're talking about AWS Lambda, milliseconds really matter. Also, because Go is fun.

### Why AWS Kool-aid?

I honestly may go back and support Azure, Google Cloud, and maybe even Bluemix later on. I think those all are developing
and Amazon has the beach head here. There's a lot of very good Amazon services - some that have been around for a while
and are battle tested and others that are new, but very exciting and powerful. The other providers will get there, but
it's going to take some time (who knows, by the time you read this, they may be there even).