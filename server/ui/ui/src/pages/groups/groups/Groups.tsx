// vendor
import React,
{
  FC,
  useCallback,
  useContext,
  useEffect,
} from 'react';
import { useMachine } from '@xstate/react';
import { useParams } from 'react-router-dom';
// environment
import { get } from 'Environment/createEnvironment';
// components
import { SectionCard } from 'Components/card/index';
import GroupsSection from './section/GroupsSection';
import Loading from '../../machine/loading/Loading';
import Error from '../../machine/error/Error';
// machine
import groupMachine from '../../machine/PageMachine';
// context
import AppContext from '../../../AppContext';
import GroupsContext from './GroupsContext';
// css
import './Groups.scss';

interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  loading?: JSX.Element;
  refetching?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}

type Data = {
  namespace?: string,
  error?: string
}


const Groups: FC = () => {
  // state
  const [state, send] = useMachine(groupMachine);
  // context
  const { user } = useContext(AppContext);
  // vars
  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';
  /**
  * Method fetches dataset data and handles state changes
  * @param {}
  * @calls {environment#get}
  * @calls {macine#send}
  * @return {void}
  */
  const fetchGroupData = useCallback(() => {
    get(`user/${user.profile.nickname}`, true).then((response: Response) => {
        return response.json();
      })
      .then((data: Data) => {
        if ((typeof data === 'string') || data.error) {
          send("ERROR", { error: data.error || data});
        } else {
          send("SUCCESS", { data: data });
        }
      })
      .catch((error: Error) => {
        const newErrorMessage = error.toString ? error.toString() : error;
        send("ERROR", {error: newErrorMessage});
      });
  }, [send, user]);



  useEffect(()=> {
    if (state.value === 'idle') {
      send('SUBMIT');
      fetchGroupData();
    }

    if (state.value === 'refetching') {
      send('SUBMIT');
      fetchGroupData();
    }
  }, [send, fetchGroupData, state.value])

  useEffect(() => {
    return () => {
      send('RESET');
    }
  }, [send]);

  const renderMap: RenderMap = {
    idle: (
      <div />
    ),
    loading: (
      <Loading />
    ),
    refetching: (
      <div />
    ),
    success: (
      <GroupsSection user={state.event.data} />
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  }


  return (
    <GroupsContext.Provider value={{ send }}>
      <div className="grid">
        <SectionCard>
          {renderMap[stateValue]}
        </SectionCard>
      </div>
    </GroupsContext.Provider>
  )
}

export default Groups;
