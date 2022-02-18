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
import { HierarchyHeader } from 'Components/header/index';
import GroupSection from './section/GroupSection';
import Loading from '../../machine/loading/Loading';
import Error from '../../machine/error/Error';
// machine
import groupMachine from '../../machine/PageMachine';
// context
import GroupContext from './GroupContext';
// css
import './Group.scss';

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

interface ParamTypes {
  groupname: string;
}



const Group: FC = () => {
  // params
  const { groupname } = useParams<ParamTypes>();
  // state
  const [state, send] = useMachine(groupMachine);
  // vars
  const stateValue: string = typeof state.value === 'string' ? state.value : 'idle';
  /**
  * Method fetches gorup data and handles state changes
  * @param {}
  * @calls {environment#get}
  * @calls {macine#send}
  * @return {void}
  */
  const fetchGroupData = useCallback(async () => {
    get(`group/${groupname}`, true).then((response: Response) => {
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
  }, [send, groupname]);



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
      <GroupSection
        group={state.event.data}
      />
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  }



  return (
    <GroupContext.Provider value={{send, groupname }}>
      <div className="grid flex flex--column">
        <h4>Group Membership</h4>
        <SectionCard>
          {renderMap[stateValue]}
        </SectionCard>
      </div>
    </GroupContext.Provider>
  );
}

export default Group;
