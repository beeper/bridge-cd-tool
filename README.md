# bridge-cd-tool

This is a simple program used to notify the Beeper API server about bridge updates.
It's probably not very interesting for others.

## Usage as a GitHub Action

```yaml
- name: Run bridge CD tool
  uses: beeper/bridge-cd-tool@main
  env:
    CI_REGISTRY: "${{ secrets.CI_REGISTRY }}"
    BEEPER_DEV_ADMIN_API_URL: "${{ secrets.BEEPER_DEV_ADMIN_API_URL }}"
    BEEPER_STAGING_ADMIN_API_URL: "${{ secrets.BEEPER_STAGING_ADMIN_API_URL }}"
    BEEPER_PROD_ADMIN_API_URL: "${{ secrets.BEEPER_PROD_ADMIN_API_URL }}"
    BEEPER_DEV_ADMIN_NIGHTLY_PASS: "${{ secrets.BEEPER_DEV_ADMIN_NIGHTLY_PASS }}"
    BEEPER_STAGING_ADMIN_NIGHTLY_PASS: "${{ secrets.BEEPER_STAGING_ADMIN_NIGHTLY_PASS }}"
    BEEPER_PROD_ADMIN_NIGHTLY_PASS: "${{ secrets.BEEPER_PROD_ADMIN_NIGHTLY_PASS }}"
```
