// vendor
import React, { FC } from 'react';
import { faDatabase } from '@fortawesome/free-solid-svg-icons';
// environment
import { get } from 'Environment/createEnvironment';
// components
import { IconButton } from 'Components/button/index';


interface ObjectStore {
  type: string;
}

interface Namespace {
  bucket_name: string;
  name: string;
  object_store: ObjectStore;
}

interface Props {
  dataset: string,
  namespaceData: Namespace;
}

const FileBrowserButton: FC<Props> = ({ dataset, namespaceData }: Props) => {

  const handleClick = () => {
    get(`namespace/${namespaceData.name}/sts`).then((response: Response) => {
        return response.json();
    }).then((data) => {
      const hash = `bucket=${namespaceData.bucket_name}&sessionToken=${data.session_token}&endpoint=${data.endpoint}&accessKeyId=${data.access_key_id}&secretAccessKey=${data.secret_access_key}&prefix=${dataset}/&delimeter=/`
      window.open(`${window.location.origin}/ui/browser/index.html#${hash}`, '__blank');
    }).catch((error: Error) => {
        console.log(error);
    });
  }

  return (
    <IconButton
      click={handleClick}
      color="white"
      icon={faDatabase}
    />
  );
}

export default FileBrowserButton;
