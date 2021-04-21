# Secret Garden

Hello, this is a tech demo of a solution to provide more control over GitHub tokens used in GitHub Actions.

Now a [similar thing is supported by GitHub](https://github.blog/changelog/2021-04-20-github-actions-control-permissions-for-github_token/), that doesn't require these hacks. Use that instead!


It relies on an embedded GitHub App, and uses the [create_access_token](https://docs.github.com/en/free-pro-team@latest/rest/reference/apps#create-an-installation-access-token-for-an-app) API call to create GitHub API tokens. These tokens can have specific permissions, and include **multiple repositories** with the same owner.
Tokens are stored as [GitHub Actions Secrets](https://docs.github.com/en/free-pro-team@latest/rest/reference/actions#secrets), so workflows can use them like [secrets.GITHUB_TOKEN](https://docs.github.com/en/free-pro-team@latest/actions/reference/authentication-in-a-workflow).

Tokens are valid for 1 hour, so it's envisioned this Action be invoked to refresh them every 15-30 minutes.

**WARNING:** Unlike the built-in `GITHUB_TOKEN`, events triggered by these tokens will trigger Actions themselves. This could cause an infinite loop, e.g. by pushing to a workflow triggered by `push`.

## Setup

1. Do not use this. Take the idea and build a better version. If you just want to try it out:
1. Fork this repo.
1. Create a [New GitHub App](https://github.com/settings/apps/new).
   Permissions must include:
     * Contents: Read
     * Metadata: Read
     * Secrets: Read & write

   The app must also have all the permissions you want to issue tokens with, see [GITHUB_TOKEN permissions](https://docs.github.com/en/free-pro-team@latest/actions/reference/authentication-in-a-workflow#permissions-for-the-github_token) for a safe selection.
1. Note the "App ID" of the app created. Store in your fork as Actions repository secret `SG_APP_ID`.
1. Create and download a private key for the app. Copy the file's contents and paste as Actions repository secret `SG_APP_PK`.
1. Install the app on your account/organization, note the "Installation ID" in the URL bar. Store as Actions repository secret `SG_INSTALL_ID`.
1. Customize the yaml in the `config/` directory.
  - `/config/secrets.yaml` defines ORG-wide tokens
  - `/config/${repo}/secrets.yaml` defines secrets for a single repository
8. Use "workflow dispatch" to test generating and storing tokens. 🤞.
9. Experiment with the tokens.


## About

This was a half-day professional development project.

The name alludes to the "seeds" from the `config/` directory growing into secrets across the organization's repositories. A small set of "gardeners" could manage access through this one repository.

Plus teh [Jerry McGuire](https://www.youtube.com/watch?v=_d_OdqErMsc) memes - it really "completes" the story of securing Actions usage.
