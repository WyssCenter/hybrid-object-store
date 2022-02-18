// vendor
import React,
  {
    FC,
    useState,
  } from 'react';
import { Link } from 'react-router-dom';
// components
import Card, { AddCard } from 'Components/card/index';
import CreateModal from '../../../shared/modal/CreateModal';
//css
import './CardList.scss';

interface Namespace {
  bucket_name: string;
  description: string;
  name: string;
  key: any;
}

interface Props {
  namespacesList: Array<Namespace>;
  send: any;
}



const CardList: FC<Props> = ({
  namespacesList,
  send,
}: Props) => {

  const [modalVisible, updateModalVisible] = useState(false)

  /**
  * Method handles close for modal
  * @param {}
  * @return {void}
  * @fires {#updateModalVisible}
  */
  const handleClose = () : void => {
    updateModalVisible(false);
  }

  return (
    <div className="CardList">

      <CreateModal
        handleClose={handleClose}
        isVisible={modalVisible}
        modalType="namespace"
        postRoute="namespace/"
        sendRefetch={() => send('REFETCH')}
      />

      <AddCard
        updateModalVisible={updateModalVisible}
        type="namespace"
      />

      {
        namespacesList.map((namespace: Namespace) => (
          <Link
            className="CardList__link"
            key={namespace.name}
            to={namespace.name}
          >
            <Card
              path={`${namespace.bucket_name}/${namespace.name}`}
              {...namespace}
            />
          </Link>
        ))
      }
    </div>
  )
}

export default CardList;
