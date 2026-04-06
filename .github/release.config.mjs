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
        preset: 'angular',
      },
    ],
    [
      '@semantic-release/release-notes-generator',
      {
        preset: 'angular',
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
