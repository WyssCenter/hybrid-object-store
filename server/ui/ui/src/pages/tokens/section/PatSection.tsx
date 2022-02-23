// vendor
import React,
{
  FC,
  useCallback,
  useEffect,
} from 'react';
import { get } from 'Environment/createEnvironment';
import { useMachine } from '@xstate/react';
// machine
import pageMachine from 'Pages/machine/PageMachine';
// components
import Error from 'Pages/machine/error/Error';
import Loading from 'Pages/machine/loading/Loading';
import { SectionCard } from 'Components/card/index';
import ListPat from './list/ListPat';


interface RenderMap {
  [key: string]: JSX.Element | undefined;
  idle?: JSX.Element;
  processing?: JSX.Element;
  refetching?: JSX.Element;
  error?: JSX.Element;
  success?: JSX.Element;
}

const PatSection: FC = () => {

  // machine
  const [state, send] = useMachine(pageMachine);

  const stateValue: string = typeof state.value === 'string'
    ? state.value
    : 'idle';

  /**
  * Method shows fetches data for pat list
  * @param {}
  * @return {void}
  * @call {machine#send}
  */
  const getList = useCallback(() => {
    get('pat/', true).then((response) => {
        return response.json();
      }).then((data) => {
        if (data && data.error) {
          send('ERROR', { error: data.error })
          return;
        }
        send('SUCCESS', { data: data });
      }).catch((error) => {

        send('ERROR', { error: error.toString() });
      });
    }, [send]);

  useEffect(() => {
    if (stateValue === 'idle') {
      send('SUBMIT');
      getList();
    }

    if (stateValue === 'refetching') {
      send('REFETCH');
      getList();
    }

  }, [getList, send, stateValue]);

  const renderMap:RenderMap = {
    idle: (
      null
    ),
    loading: (
      <Loading />
    ),
    refetching: (
      <div />
    ),
    success: (
      <SectionCard>
        <div>
          <ListPat
            list={state.event.data}
            send={send}
          />
        </div>
      </SectionCard>
    ),
    error: (
      <Error errorMessage={state.event.error} />
    )
  };

  return (
    renderMap[stateValue]
  )
}

export default PatSection;
