// vendor
import React, { FC, useCallback, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import ReactTooltip from 'react-tooltip';
// environemnt
import { put, del } from 'Environment/createEnvironment';
// components
import Delete from 'Shared/delete/Delete';
import {SectionCard} from 'Components/card/index';
import {WarningButton} from 'Components/button/index';
import DeleteModal from './modal/DeleteModal';
// assets
import './DeleteDataset.scss';

// types
import {
  Dataset,
} from '../section/DatasetSectionTypes';

interface Props {
  send: any;
  route: string;
  redirectRoute: string;
  pendingDelete: boolean;
}

interface ParamTypes {
  datasetName: string;
  namespace: string;
}

let deleteTime = 0;

const baseHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Content-Type": 'application/json',
  "Origin": window.location.origin
}

const baseUrl = `${window.location.protocol}//${window.location.hostname}/core/v1`;

fetch(`${baseUrl}/discover`, { headers: baseHeaders, method: 'GET' })
.then(response => response.json())
.then(data => deleteTime = Math.round(data.delete_delay_minutes / 60));

const DeleteDataset: FC<Props> = ({ send, redirectRoute, route, pendingDelete }: Props) => {
  // state
  const [hasError, setHasError] = useState(null);
  const [ isModalVisible, setIsModalVisible ] = useState(false);
  // params
  const { datasetName, namespace } = useParams<ParamTypes>();
  // click functions

  const handleDeleteClick = () => {
    del(route).then((response: Response) => {
      window.location.pathname = redirectRoute;
    }).catch((error: Error) => {
      setHasError(`Error: Dataset may have not been deleted correctly, refresh to confirm`);
    });
  }

  return (
    <section>
      <h3 className="DatasetSection__h3">Delete Dataset</h3>

      <DeleteModal
        deleteHours={deleteTime}
        datasetName={datasetName}
        handleDeleteClick={handleDeleteClick}
        hideModal={() => setIsModalVisible(false)}
        isVisible={isModalVisible}
      />
      <SectionCard>
        <div className="DeleteDataset flex align-items--center">
          <p>{
              `Deleting a dataset will schedule it for removal. Once a dataset is marked for deletion, it will be deleted after ${deleteTime} hours. If it is not restored by an administrator before this time, all data will be permanently removed!`
          }</p>
          <div
            className="DeleteDataset__buttons relative"
          >
          <div
            data-tip="Dataset deletion is already pending."
            role="presentation"
            data-tip-disable={!pendingDelete}
          >
            <WarningButton
              click={() => setIsModalVisible(true)}
              text="Delete Dataset"
              disabled={pendingDelete}
            />
            {
              pendingDelete && (<ReactTooltip
                place="bottom"
                effect="solid"
              />)
            }
          </div>
          </div>
        </div>
      </ SectionCard>
    </section>
  );
}


export default DeleteDataset;
