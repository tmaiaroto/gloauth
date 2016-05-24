# GLOAuth

Go Lambda OAuth ("glow auth") is a serverless authentication system relying on AWS (Lambda, DynamoDB, Cognito, SES, API Gateway).

It's goal is to build a simple, yet flexible, zero maintenance authentication system for all of your application needs. 
Why another one? Well, there's a lot of good auth packages for Node.js (Passport for example) but I haven't find any for
Golang nor do I find many that leverage AWS Lambda (there is LambdAuth, again in Node.js). At least not at the time I needed
this. If I did, I wouldn't be building this, trust me.

### Setup

GLOAuth is using Apex and Terraform to orchestrate a few things. You'll of course need an AWS account. Ensure it's configured
and you have your `~/.aws/config` information set OR have environment variables for your credentials. You'll need to provide
your `account_id` when you run Terraform. Apex wraps Terraform but doesn't prompt you for variables. You can pass your account
id in -var and set any other variable you need this way too. You could also alter the `variables.tf` in the `dev` and `prod`
infrastructure directories as needed. It's a bit cleaner and safer to rely on your AWS config and to pass your account id
when running the infra command.

```
apex deploy
apex infra get
apex infra apply -var 'aws_account_id=12345'
```

...TBD...

### Why Go?

Well, LambdAuth wasn't exactly turn-key. The setup process wasn't working out for me and I quickly realized that not only
was it not going to be as "easy" as I thought, but it was also done with Node.js so it's slower. Not a lot of course, but
when we're talking about AWS Lambda, milliseconds really matter. Also, because Go is fun.

### Why AWS Kool-aid?

I honestly may go back and support Azure, Google Cloud, and maybe even Bluemix later on. I think those all are developing
and Amazon has the beach head here. There's a lot of very good Amazon services - some that have been around for a while
and are battle tested and others that are new, but very exciting and powerful. The other providers will get there, but
it's going to take some time (who knows, by the time you read this, they may be there even).