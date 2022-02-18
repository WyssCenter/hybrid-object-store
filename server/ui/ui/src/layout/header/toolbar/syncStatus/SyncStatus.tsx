// vendor
import React, { FC } from 'react';
// css
import { ReactComponent as Disabled } from 'Images/icons/icon-sync-disabled.svg';
import { ReactComponent as Simplex } from 'Images/icons/icon-sync-simplex.svg';
import { ReactComponent as Duplex } from 'Images/icons/icon-sync-duplex.svg';
import './SyncStatus.scss';



interface Props {
  sync: any[];
}


const SyncButton:FC<Props> = ({
  sync,
}: Props) => {

  if (!sync) {
    return <div />;
  }

  // namespace styling section
  const syncEnabled = sync[0] && sync[0].data.length > 0;
  let namespaceSyncValue = 0;
  let namespaceBadgeValue = 'Disabled';

  if (syncEnabled) {
    if (sync[0].data[0].sync_type === 'simplex') {
      namespaceSyncValue = 1;
      namespaceBadgeValue = '1-way';
    } else if (sync[0].data[0].sync_type === 'duplex') {
      namespaceSyncValue = 2;
      namespaceBadgeValue = '2-way';
    }
  }

  // dataset styling section
  const isDataset = sync.length > 1;
  let datasetSyncValue = 0;
  let datasetBadgeValue = 'Disabled';

  if (isDataset && sync && sync[1] && sync[1].data.syncEnabled) {
    if (sync[1].data.syncType === 'simplex') {
      datasetSyncValue = 1;
      datasetBadgeValue = '1-way';
    } else if (sync[1].data.syncType === 'duplex') {
      datasetSyncValue = 2;
      datasetBadgeValue = '2-way';
    }
  }

  return (
    <div className="SyncStatus">
      <div className="SyncStatus__text">
          <span className="SyncStatus__label">
            Namespace Sync:
          </span>
          <div className="SyncStatus__icon">
            {(namespaceSyncValue === 0) && <Disabled />}
            {(namespaceSyncValue === 1) && <Simplex />}
            {(namespaceSyncValue === 2) && <Duplex />}
          </div>
          <span className="SyncStatus__value">{namespaceBadgeValue}</span>
      </div>
      {
        isDataset && (
      <div className="SyncStatus__text">
          <span className="SyncStatus__label">
            Dataset Sync:
          </span>
          <div className="SyncStatus__icon">
            {(datasetSyncValue === 0) && <Disabled />}
            {(datasetSyncValue === 1) && <Simplex />}
            {(datasetSyncValue === 2) && <Duplex />}
          </div>
          <span className="SyncStatus__value">{datasetBadgeValue}</span>
        </div>
        )
      }
    </div>
  );
}

export default SyncButton;
