package config

const exampleConfig = "" +
	`local User(name, email) = { name: name, email: email };

local user = {
  personal: User('jdoe', 'jdoe@users.noreply.github.com'), // Example user
  work: User('John Doe', 'john.doe@corp.com'), // Example user
};

// This is just an example.
// Feel free to delete it.
local example = {
  'github.com': {  // Public GitHub
    default: user.personal,
    overrides: [
      {  // Your employer's organization at github.com
        owner: 'corp',
        user: user.work,
      },
    ],
  },
  'github.corp.com': {  // Your employer's GitHub Enterprise
    default: user.work,
  },
}

// This is your real config.
local config = {}

config
`
