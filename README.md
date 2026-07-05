# Meshbox Studio

Meshbox Studio is a lightweight, self-hosted application designed for 3D printing enthusiasts to archive their projects, document their progress, and keep their digital maker space organized.

Tired of losing your perfect slicer settings or forgetting why a print failed three months ago? Ever need to track down a file you printed last year for Christmas?
Meshbox Studio acts as your personal 3D printing journal and file vault.

## Disclaimer

> [!IMPORTANT]
> This project is still in early development and not recommended for general use.
> We are currently focusing on architectural foundations and welcome early feedback or contributions to help shape the roadmap.

## Features

- Import 3D models from Printables .zip exports (drag & drop)
- Auto-extract metadata: title, designer, tags, license, print stats, and cover image
- Archive, search, and organize projects with edit support
- Track print iterations with notes and outcomes

## Repository Architecture

- `backend/`: Go API + static file server (embeds frontend build output)
- `frontend/`: Nuxt 4 SPA (Nuxt UI + Tailwind CSS)

For production, the frontend is built into `backend/internal/webui/dist` and served by the Go backend as a single self-hosted binary.

## License

This project is licensed under the GNU Affero General Public License v3.0.

See the full license text in [LICENSE](LICENSE).
