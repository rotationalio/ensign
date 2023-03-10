import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
import { FormikProvider, useFormik } from 'formik';
import styled from 'styled-components';

type GenerateAPIKeyModalProps = {
  open: boolean;
  setOpenGenerateAPIKeyModal: React.Dispatch<React.SetStateAction<boolean>>;
};

function GenerateAPIKeyModal({ open, setOpenGenerateAPIKeyModal }: GenerateAPIKeyModalProps) {
  const formik = useFormik({
    initialValues: {},
    onSubmit: () => {},
  });

  return (
    <Modal
      open={open}
      title={<h1>Generate Your API Key</h1>}
      size="medium"
      onClose={() => setOpenGenerateAPIKeyModal(false)}
      containerClassName="overflow-scroll h-[80vh] lg:h-[90vh]"
    >
      <FormikProvider value={formik}>
        <div>
          <p className="mb-5">
            Name your key and customize permissions. Or stick with the default.
          </p>
          <form className="space-y-6">
            <fieldset>
              <h2 className="mb-3 font-semibold">Key Name</h2>
              <TextField placeholder="default" fullWidth />
            </fieldset>
            <fieldset>
              <h2 className="mb-3 font-semibold">Permissions</h2>
              <div className="space-y-8">
                <Box>
                  <h2 className="mb-1 font-semibold">Full Access</h2>
                  <StyledFieldset>
                    <Checkbox {...formik.getFieldProps('all')}>
                      Full Access (default) - Publish to topic, Subscribe to topic, Create Topic,
                      Read Topic, Delete Topic, Destroy Topic.
                    </Checkbox>
                  </StyledFieldset>
                </Box>
                <Box>
                  <h2 className="mb-1 font-semibold">Custom Access</h2>
                  <StyledFieldset>
                    <Checkbox {...formik.getFieldProps('custom')}>
                      Check to grant access for each action.
                    </Checkbox>
                  </StyledFieldset>
                  <div className="mt-5 ml-5 w-full space-y-1 md:ml-10 md:w-1/2">
                    <StyledFieldset>
                      <Checkbox>Publish to topic</Checkbox>
                    </StyledFieldset>
                    <StyledFieldset>
                      <Checkbox>Subscribe to topic</Checkbox>
                    </StyledFieldset>
                    <StyledFieldset>
                      <Checkbox>Create topic</Checkbox>
                    </StyledFieldset>
                    <StyledFieldset>
                      <Checkbox>Read topic</Checkbox>
                    </StyledFieldset>
                    <StyledFieldset>
                      <Checkbox>Delete topic</Checkbox>
                    </StyledFieldset>
                    <StyledFieldset>
                      <Checkbox>Destroy topic</Checkbox>
                    </StyledFieldset>
                  </div>
                </Box>
              </div>
            </fieldset>
            <div className="flex items-center justify-between">
              <p className="pl-5">
                Generate API Key for <span className="font-semibold">[insert project name]</span>{' '}
                project.
              </p>
              <Button className="bg-[#6DD19C] px-6 py-3 font-semibold">Generate API Key</Button>
            </div>
          </form>
        </div>
      </FormikProvider>
    </Modal>
  );
}

export default GenerateAPIKeyModal;

const Box = styled.div`
  padding: 16px 20px;
  border: 1px solid gray;
  border-radius: 4px;
`;

const StyledFieldset = styled.fieldset`
  label {
    flex-direction: row-reverse !important;
    justify-content: space-between;
    font-size: 14px;
  }
  label svg {
    min-width: 25px;
  }
`;
