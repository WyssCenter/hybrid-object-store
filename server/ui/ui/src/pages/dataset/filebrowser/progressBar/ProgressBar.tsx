// vendor
import React, { FC, useState, useEffect } from 'react';
// css
import './ProgressBar.scss';


interface Props {
  uploadData: any;
}


const StaticButtons:FC<Props> = ({
  uploadData,
}: Props) => {

  const root = document.documentElement;
  const primaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#957299';

  const pct = (uploadData.finishedSize / uploadData.totalSize * 100)


  useEffect(() => {
    const advanceProgress = ( amount: any ) => {
      amount = amount || 0;
      if(document.getElementById( 'progressCount' )) {
        document.getElementById( 'progressCount' ).innerHTML = amount.toFixed(2);
      }
    }

    advanceProgress(pct);

  }, [uploadData])

  const style = { "--pct": `${pct}%` } as React.CSSProperties;

  return (
    <div className="ProgressBar">
      <div id="progressBar" style={style}>
        <span id="progressText">
          <span className="ProgressBar__count">{`Uploading ${uploadData.totalFiles} File${uploadData.totalFiles.length >1 ? 's' : ''}`}</span>
          <span className="ProgressBar__pct">
            {(Number(pct) >= 100 )  &&
              'Finalizing Upload...'
            }
            {(Number(pct) < 100) &&
              (<><span id="progressCount">0</span>% Complete</>)
            }
          </span>
        </span>
      </div>
    </div>
  );
}

export default StaticButtons;
