// vendor
import React,
{
  FC,
  useCallback,
  useContext,
  useState,
} from 'react';
import { Link } from 'react-router-dom';
import { faChevronRight, faChevronDown } from '@fortawesome/free-solid-svg-icons';
// environment
import { get } from 'Environment/createEnvironment';
// app context
import AppContext from 'Src/AppContext';
// componets
import Delete from 'Shared/delete/Delete';
import { IconButton } from 'Components/button/index';
import NamespaceDetails from './details/NamespaceDetails';
// context
import NamespaceListingContext from '../../NamespaceListingContext';
// css
import { ReactComponent as Disabled } from 'Images/icons/icon-sync-disabled.svg';
import { ReactComponent as Simplex } from 'Images/icons/icon-sync-simplex.svg';
import { ReactComponent as Duplex } from 'Images/icons/icon-sync-duplex.svg';
import './NamespaceRow.scss';


interface ObjectStore {
  type: string;
}


interface Namespace {
  bucket_name: string;
  description: string;
  directory: string;
  name: string;
  object_store: ObjectStore;
}

interface Props {
  namespace: Namespace;
}


const NamespaceRow: FC<Props> = ({ namespace }: Props ) => {
  // context
  const { user } = useContext(AppContext);
  // state
  const [isExpanded, setIsExpanded] = useState(false);
  const [syncStatus, setSyncStatus] = useState('...');

  get(`namespace/${namespace.name}/sync`).then((response: Response) => {
      return response.json();
    })
    .then((data) => {
      if(data && !data[0]) {
        setSyncStatus('disabled');
      } else {
        setSyncStatus(data[0].sync_type);
      }
    })

  /**
  * Method toggles detail view
  * @param {}
  * @return {void}
  * @calls {state#setIsExpanded}
  */
  const toggleView = useCallback(() => {
    setIsExpanded(!isExpanded);
  }, [isExpanded])

  const icon = isExpanded ? faChevronDown : faChevronRight;

  let sanitizedSyncName = '...';
  if (syncStatus === 'simplex') {
    sanitizedSyncName = '1-Way';
  } else if (syncStatus === 'duplex') {
    sanitizedSyncName = '2-Way';
  } else if (syncStatus === 'disabled') {
    sanitizedSyncName = 'Disabled';
  }

  return (
    <>
      <tr>
        <td className="NamespaceRow__td--empty">
          <IconButton
            click={toggleView}
            color="primary"
            backgroundColor="transparent"
            icon={icon}
          />
        </td>
        <td>
          <Link to={`/${namespace.name}`}>
            {namespace.name}
          </Link>
        </td>
        <td>{namespace.description}</td>
        <td>{namespace.object_store.type}</td>
        <td>{namespace.bucket_name}</td>
        <td className={`NamespaceRow__sync`}>
          {(syncStatus === 'disabled') && <Disabled />}
          {(syncStatus === 'simplex') && <Simplex />}
          {(syncStatus === 'duplex') && <Duplex />}
          {sanitizedSyncName}
        </td>

      </tr>

      <NamespaceDetails
        namespaceName={namespace.name}
        isExpanded={isExpanded}
      />

    </>
  );
}


export default NamespaceRow;
