// vendor
import React,
{
  FC,
  useCallback,
  useRef,
  useState
} from 'react';
// components
import PrimaryButton from 'Components/button/index';
import { CopyText } from 'Components/form/text/index';
// css
import './NewPat.scss';

interface PAT {
  token: string;
  id: number;
}

interface Props {
  dismissPat: () => void;
  pat: PAT;
}

const NewPat: FC<Props> = ({ dismissPat, pat }:Props) => {
  // ref
  const tooltipRef = useRef(null);
  // state
  const [isTooltipVisible, setTooltipVisible] = useState(false);

  /**
  * Method updates tooltip visibility to true
  * @param {}
  * @return {void}
  * @call {state#setTooltipVisible}
  */
  const showTooltipVisible = useCallback(() => {
    setTooltipVisible(true);
  }, []);

  /**
  * Method updates tooltip visibility to false
  * @param {}
  * @return {void}
  * @call {state#setTooltipVisible}
  */
  const hideTooltipVisible = useCallback(() => {
    setTooltipVisible(false);
  }, []);

  if (pat == null) {
    return <div className="CreatePat__filler" />;
  }

  return (
    <div className="NewPat__Container">
      <div className="NewPat">
        <p>Copy personal access token now. This token will not viewable again!</p>
        <CopyText text={pat.token} />
      </div>
          <div ref={tooltipRef} className="NewPat__Buttons">
            <br />
            <PrimaryButton
              click={dismissPat}
              text="Close"
            />
          </div>
    </div>
  );
}



export default NewPat;
