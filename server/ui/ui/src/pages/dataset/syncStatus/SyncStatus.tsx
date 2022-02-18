// vendor
import React, { FC } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSyncAlt } from '@fortawesome/free-solid-svg-icons';
// css
import './SyncStatus.scss';

interface Dataset {
  sync_enabled: false;
}


interface Props {
  dataset: Dataset;
}


const SyncButton:FC<Props> = ({
  dataset
}: Props) => {
  const isEnabled = dataset && (dataset.sync_enabled)
  const colorButton = isEnabled ? '#39b983' : '#9b9c9e';
  const syncText = isEnabled ? 'Sync Enabled' : 'Sync Disabled';


  return (
    <div className="SyncStatus">
      <b>{syncText}</b>
      <FontAwesomeIcon
        color={colorButton}
        icon={faSyncAlt}
        size="lg"
      />
    </div>
  );
}

export default SyncButton;
