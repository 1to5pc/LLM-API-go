# Golang LLM Accessor with OpenRouter API
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/1to5pc/LLM-API-go/app-build.yml?style=for-the-badge)

A Golang project that allows users to access large language models (LLMs) via the OpenRouter API. The application stores user inputs, retains context across LLM switches, and provides a sleek, user-friendly interface for a seamless experience.

## Features

- **Access LLMs via OpenRouter API**: Easily interact with LLMs through the OpenRouter API.
- **Save User Inputs**: User inputs are saved, providing context for continuous conversations.
- **Switch Between LLMs**: Effortlessly switch between different LLMs while retaining context across sessions.
- **Beautiful UI**: An intuitive, visually pleasing UI to enhance user experience.
- **Context Retention**: Keep context intact when switching between different LLMs for smooth and coherent interactions.

## Installion
[![GitHub Release](https://img.shields.io/github/v/release/1to5pc/LLM-API-go?display_name=tag&style=for-the-badge)](https://github.com/1to5pc/LLM-API-go/releases/latest)

1. Download the applicable version of the program for your OS using the above link.
2. Execute the program and enjoy!

> [!IMPORTANT]
> For the program to function as intended an OpenRouter API key needs to provided. See below.

## Adding an OpenRouter API Key

To access the OpenRouter API, you'll need an API key. Follow these steps to get and configure the key for your environment:

### 1. Get an API Key from OpenRouter

1. Go to [OpenRouter API Keys](https://openrouter.ai/settings/keys).
2. Sign in or create an account if you haven't already.
3. Generate a new API key, and copy it.

### 2. Export the API Key in Your Operating System

Once you have your API key, follow the instructions below based on your OS to export the API key as an environment variable.

#### For Linux / macOS

1. Open your terminal.
2. Use the following command to export the API key temporarily:

   ```bash
   export OPENROUTER_API_KEY=your-api-key-here
   ```

   To make this change permanent across terminal sessions, add the export command to your shell profile:

   - For **Bash** users, add the line to your `~/.bashrc`:

     ```bash
     echo 'export OPENROUTER_API_KEY=your-api-key-here' >> ~/.bashrc
     source ~/.bashrc
     ```

   - For **Zsh** users, add the line to your `~/.zshrc`:

     ```bash
     echo 'export OPENROUTER_API_KEY=your-api-key-here' >> ~/.zshrc
     source ~/.zshrc
     ```

3. To verify that the API key was set successfully, run:

   ```bash
   echo $OPENROUTER_API_KEY
   ```

#### For Windows

1. Open **Command Prompt** (or **PowerShell**).
2. Set the environment variable using the following command:

   For **Command Prompt**:

   ```cmd
   setx OPENROUTER_API_KEY "your-api-key-here"
   ```

   For **PowerShell**:

   ```powershell
   [System.Environment]::SetEnvironmentVariable('OPENROUTER_API_KEY', 'your-api-key-here', [System.EnvironmentVariableTarget]::User)
   ```

3. Restart your terminal for the changes to take effect. You can verify it by running:

   ```cmd
   echo %OPENROUTER_API_KEY%
   ```

After setting the environment variable on your respective system, you can proceed to configure the application to use the OpenRouter API key.

### 3. Add the API Key to the Project

To securely store the key in your project, create a `.env` file in the root directory of the project and add the following line:

```env
OPENROUTER_API_KEY=your-api-key-here
```

This ensures the application can retrieve the key from the environment variable or `.env` file during runtime.

---

## Development
Clone the repository and get started right away.

```bash
git clone https://github.com/1to5pc/LLM-API-go.git
```
