local User(name, email) = { name: name, email: email };

local user = {
  "work": User("John Doe", "john.doe@corp.com")
};

{
  "github.corp.com": {
    "default": user.work
  }
}
