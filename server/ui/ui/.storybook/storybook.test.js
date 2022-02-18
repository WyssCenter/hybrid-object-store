import path from 'path';

import initStoryshots, { Stories2SnapsConverter } from '@storybook/addon-storyshots';

jest.mock('react-dom');

// function to customize the snapshot location
const getMatchOptions = ({ context: { fileName } }) => {
  // Generates a custom path based on the file name and the custom directory.
  const snapshotPath = path.join(path.dirname(fileName), '__snapshots__');
  return { customSnapshotsDir: snapshotPath };
};

const user = JSON.stringify({ id_token: 'id_token' });
localStorage.setItem(`oidc.user:http://localhost/auth/v1/.well-known/openid-configuration:HossServer`, user);

const modalElement = window.document.createElement('div');
modalElement.setAttribute('id', 'modal');

document.body.appendChild(modalElement)


const beforeScreenshot = (page, { context: { kind, story }, url }) => {
  page.localStorage = new LocalStorageMock();
  const modalElement = global.document.createElement('div');
  modalElement.setAttribute('id', 'modal');
  const currentModalElement = global.document.getElementById('modal');

  global.document.body.removeChild(currentModalElement);

  global.document.body.appendChild(modalElement);



  return new Promise((resolve) =>
    setTimeout(() => {
      resolve();
    }, 600)
  );
};


initStoryshots({
  stories2snapsConverter: new Stories2SnapsConverter({
    snapshotExtension: '.storypuke',
    storiesExtensions: ['stories.tsx'],
    beforeScreenshot,
  }),
});
