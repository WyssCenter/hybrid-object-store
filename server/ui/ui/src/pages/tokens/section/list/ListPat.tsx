// vendor
import React, { FC, useState} from 'react';
import { faPlus } from '@fortawesome/free-solid-svg-icons';
// components
import { FlatIconTextButton } from 'Components/button/index';
import PatItem from './item/PatItem';
import CreatePatModal from './create/CreatePatModal';

// css
import './ListPat.scss';

interface PAT {
  description: string;
  id: number;
}
interface Props {
  list: Array<PAT>;
  send: any;
}

const ListPat: FC<Props> = ({ list, send }: Props) => {
  // state
  const [ isModalVisible, setIsModalVisible] = useState(false);


  return(
    <>
      <CreatePatModal
        handleClose={() =>setIsModalVisible(false)}
        isVisible={isModalVisible}
        send={send}
      />
      <table className="ListPat__table column-1-span-6">
        <thead>
        <tr>
          <th className="column-1-span-4">Description</th>
          <th>Revoke</th>
        </tr>
        </thead>
        <tbody>
          <tr>
            <td
              className="td--center  td--actions"
              colSpan={2}
            >
              <FlatIconTextButton
                click={() => setIsModalVisible(true)}
                color="primary"
                icon={faPlus}
                text="Create a Personal Access Token"
              />
            </td>
          </tr>
          { list && list.map((pat: PAT) => (
              <PatItem
                key={pat.id}
                pat={pat}
                send={send}
              />
            ))
          }
          { list && list.length === 0 && (
            <td className="ListPat__td--empty" colSpan={2}>
              <p>Currently you have no personal access tokens.</p>
            </td>
          )


          }
        </tbody>
      </table>
    </>
  )
}


export default ListPat;
