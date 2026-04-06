/**
 * @type {import('semantic-release').GlobalConfig}
 */

const isRelease = process.env.IS_RELEASE === 'true';

export default {
  debug: true,
  branches: isRelease ? ['main'] : [{ name: 'main', prerelease: 'prerelease' }],
  tagFormat: '${version}',
  plugins: [
    [
      '@semantic-release/commit-analyzer',
      {
        preset: 'conventionalcommits',
      },
    ],
    [
      '@semantic-release/release-notes-generator',
      {
        preset: 'conventionalcommits',
      },
    ],
    // [
    //   '@semantic-release/github',
    //   {
    //     discussionCategoryName: 'Announcements',
    //     releasedLabels: false,
    //     draftRelease: isRelease,
    //   },
    // ],
  ],
};
