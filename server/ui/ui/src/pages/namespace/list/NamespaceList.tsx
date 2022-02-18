// vendor
import React, {
  FC,
  useCallback,
  useContext,
  useState,
} from 'react';
import {
  useParams,
} from "react-router-dom";
import { faPlus } from '@fortawesome/free-solid-svg-icons';
// app context
import AppContext from 'Src/AppContext';
// components
import { SectionCard } from 'Components/card/index';
import CreateModal from 'Src/shared/modal/CreateModal';
import { FlatIconTextButton } from 'Components/button/index';
import DatasetItem from './item/DatasetItem';
import SectionFilter from 'Pages/shared/filters/SectionFilter';
// context
import NamespaceContext from '../NamespaceContext'
//css
import './NamespaceList.scss'

interface ParamTypes {
  namespace: string;
}

interface Dataset {
  created: string;
  description: string;
  directory: string;
  name: string;
  path: string;
  sync: string;
  showDeleteBadge: boolean;
}

interface ObjectStore {
  type: string;
}

interface Namespace {
  bucket_name: string;
  name: string;
  object_store: ObjectStore;
}


interface Props {
  datasets: Array<Dataset>;
  namespaceData: Namespace;
}

const NamespaceListing: FC<Props> = ({ datasets, namespaceData }: Props) => {
  // context
  const { user } = useContext(AppContext);
  const { send } = useContext(NamespaceContext);
  // params
  const { namespace } = useParams<ParamTypes>();
  // state
  const [isModalVisible, updateModalVisible] = useState(false);
  const [visibleDatasets, setVisibleDatasets] = useState(datasets);



  // functions
  /**
  * Method handles close for modal
  * @param {}
  * @return {void}
  * @fires {#updateModalVisible}
  */
  const handleClose = useCallback(() : void => {
    updateModalVisible(false);
  }, [updateModalVisible]);

  /**
  * Method sends refetch method
  * @param {}
  * @return {void}
  * @fires {#send}
  */
  const sendRefetch = useCallback(() => {
    send('REFETCH');
  }, [send]);

  const permissions = user.profile.role !== 'user';

  return (
    <div className="NamespaceList">

      { permissions &&
        <CreateModal
          handleClose={handleClose}
          isVisible={isModalVisible}
          modalType="Dataset"
          postRoute={`namespace/${namespace}/dataset/`}
          sendRefetch={sendRefetch}
        />
      }

      <SectionCard>
        <div className="Namespace__actions">
          <SectionFilter
            list={datasets}
            formattedSection="Dataset"
            updateList={setVisibleDatasets}
          />
        </div>

        <table>
          <thead>
            <tr>
              <th className="NamespaceList__th--empty">
              </th>
              <th>
                Name
              </th>

              <th>
                Description
              </th>

              <th>
                Created
              </th>

              <th>
                Sync Status
              </th>
            </tr>
          </thead>
          <tbody>
              { permissions &&
                <tr>
                  <td className="td--center td--actions" colSpan={6}>
                    <FlatIconTextButton
                      click={() => updateModalVisible(true)}
                      color="primary"
                      icon={faPlus}
                      text="Create Dataset"
                    />
                  </td>
                </tr>
              }
              { visibleDatasets && visibleDatasets.map((dataset) => (
                <DatasetItem
                  dataset={dataset}
                  namespaceData={namespaceData}
                  key={dataset.name}
                />
              ))}
          </tbody>
        </table>
        {
          !visibleDatasets.length && (
            <div className ="NamespaceList__empty">
              You do not have access to any datasets.
            </div>
          )
        }
      </SectionCard>

    </div>
  )
}

export default NamespaceListing;
