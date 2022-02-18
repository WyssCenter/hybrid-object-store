// vendor
import React, {
  FC,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react';
// app context
import AppContext from 'Src/AppContext';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
// components
import { SectionCard } from 'Components/card/index';
import { FlatIconTextButton } from 'Components/button/index';
import CreateModal from 'Src/shared/modal/CreateModal';
import NamespaceRow from './row/NamespaceRow';
import SectionFilter from 'Pages/shared/filters/SectionFilter';
// context
import NamespaceListingContext from '../NamespaceListingContext';
// environment
import { get } from 'Environment/createEnvironment';
// css
import './NamespacesTable.scss';

interface ObjectStore {
  type: string;
}

interface Namespace {
  bucket_name: string;
  created: string;
  description: string;
  directory: string;
  name: string;
  object_store: ObjectStore;
}

interface Props {
  namespaces: Array<Namespace>;
}

const NamespaceTable: FC<Props> = ({ namespaces }: Props) => {
  // context
  const { user } = useContext(AppContext);
  const { send } = useContext(NamespaceListingContext);
  // state
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [visibleNamespaces, setVisibleNamespaces] = useState(namespaces);
  const [objectStoreList, setObjectStoreList] = useState([]);
  const [errorMessage, setErrorMessage] = useState(null);


  useEffect(() => {
    get('object_store/').then((response) => {
      return response.json();
    }).then((data) => {
      setObjectStoreList(data);
    }).catch((error) => {
      setErrorMessage(error.toString());
    });
  }, [])

  // functions

  const handleOpen = useCallback(() : void => {
    setIsModalVisible(true);
  }, [setIsModalVisible]);
  /**
  * Method handles close for modal
  * @param {}
  * @return {void}
  * @fires {#updateModalVisible}
  */
  const handleClose = useCallback(() : void => {
    setIsModalVisible(false);
  }, [setIsModalVisible]);

  /**
  * Method sends refetch method
  * @param {}
  * @return {void}
  * @fires {#send}
  */
  const sendRefetch = useCallback(() => {
    send('REFETCH');
  }, [send]);

  if (errorMessage) {
    return (
      <div className="NamespacesTable">
        <h5>There was a problem loading this page</h5>
        <p>errorMessage</p>
      </div>
    );
  }

  return (
    <div className="NamespacesTable">
      { (user.profile.role === 'admin') &&
        <CreateModal
          handleClose={handleClose}
          isVisible={isModalVisible}
          modalType="namespace"
          objectStoreList={objectStoreList}
          postRoute="namespace/"
          sendRefetch={sendRefetch}
        />
      }
      <SectionCard>
        <div className="NamespacesTable__actions">
          <SectionFilter

            list={namespaces}
            formattedSection="namespace"

            updateList={setVisibleNamespaces}
          />
        </div>

        <table className="NamespacesTable__table">
          <thead>
            <tr>
              <th className="NamespacesTable__th--empty">
              </th>
              <th>
                Name
              </th>

              <th>
                Description
              </th>

              <th>
                Object Store Type
              </th>

              <th>
                Bucket Name
              </th>
              <th>
                Sync Status
              </th>
            </tr>
          </thead>
          <tbody>
            { (user.profile.role === 'admin') &&
              <tr>
                <td className="NamespaceListing__td--center td--center td--actions" colSpan={7}>

                  <FlatIconTextButton
                    click={handleOpen}
                    color="primary"
                    icon={faPlus}
                    text="Create Namespace"
                  />
                </td>
              </tr>
            }
            { visibleNamespaces && visibleNamespaces.map((namespace) => (
              <NamespaceRow
                namespace={namespace}
                key={namespace.name}
              />
            ))}
          </tbody>
        </table>
      </SectionCard>

    </div>
  )
}

export default NamespaceTable;
