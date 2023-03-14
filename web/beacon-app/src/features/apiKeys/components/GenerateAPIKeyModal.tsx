/* eslint-disable prettier/prettier */
import { Button, Checkbox, Modal, TextField } from '@rotational/beacon-core';
import { Form, FormikProvider, useFormik } from 'formik';
import { useState } from 'react';
import styled from 'styled-components';

import { Close as CloseIcon } from '@/components/icons/close';
import { Toast } from '@/components/ui/Toast';
import { useCreateProjectAPIKey } from '@/features/apiKeys/hooks/useCreateApiKey';
import { APIKeyDTO, NewAPIKey } from '@/features/apiKeys/types/createApiKeyService';
import { useFetchPermissions } from '@/hooks/useFetchPermissions';
import { useOrgStore } from '@/store';
type GenerateAPIKeyModalProps = {
  open: boolean;
  onSetKey: React.Dispatch<React.SetStateAction<any>>;
  onClose: () => void;
  setOpenAPIKeyDataModal: () => void;
};
// if selected, then all permissions are selected
// if not selected, then all permissions are not selected

function GenerateAPIKeyModal({
  open,
  setOpenAPIKeyDataModal,
  onSetKey,
  onClose,
}: GenerateAPIKeyModalProps) {
  const [fullSelected, setFullSelected] = useState(true);
  const [customSelected, setCustomSelected] = useState(false);
  const org = useOrgStore.getState() as any;
  const { permissions } = useFetchPermissions();
  const { createProjectNewKey, key, wasKeyCreated, isCreatingKey, hasKeyFailed, error } =
    useCreateProjectAPIKey();
  const handleCreateKey = ({ name, permissions }: any) => {
    const payload = {
      projectID: org.projectID,
      name,
      permissions,
    } satisfies APIKeyDTO;

    createProjectNewKey(payload);

    // TODO: create handle error abstraction
  };
  if (wasKeyCreated) {
    onSetKey(key);
    onClose();
    setOpenAPIKeyDataModal();
  }

  if (hasKeyFailed || error) {
    <Toast
      isOpen={hasKeyFailed}
      variant="danger"
      description={(error as any)?.response?.data?.error || 'Something went wrong'}
    />;
  }
  const formik = useFormik<NewAPIKey>({
    initialValues: {
      name: '',
      permissions: [''],
    },
    onSubmit: (values) => {
      console.log('values', values);
      handleCreateKey(values as APIKeyDTO);
    },
    // validationSchema: NewAPIKEYSchema,
  });

  const { values, setFieldValue } = formik;

  return (
    <Modal
      open={open}
      title={<h1>Generate Your API Key</h1>}
      size="medium"
      containerClassName="overflow-hidden h-[90vh] "
    >
      <>
        <Button
          variant="ghost"
          className="bg-transparent absolute -right-10 top-5 border-none border-none p-2 p-2"
        >
          <CloseIcon onClick={close} />
        </Button>
        <FormikProvider value={formik}>
          <div>
            <p className="mb-5">
              Name your key and customize permissions. Or stick with the default.
            </p>
            <Form className="space-y-6">
              <fieldset>
                <h2 className="mb-3 font-semibold">Key Name</h2>
                <TextField placeholder="default" fullWidth {...formik.getFieldProps('name')} />
              </fieldset>
              <fieldset>
                <h2 className="mb-3 font-semibold">Permissions</h2>
                <div className="space-y-8">
                  <Box>
                    <h2 className="mb-1 font-semibold">Full Access</h2>
                    <StyledFieldset>
                      <Checkbox
                        {...formik.getFieldProps('full')}
                        onChange={(isSelected) => {
                          setFullSelected(!!isSelected);
                          setCustomSelected(false);
                          // reset permissions
                          setFieldValue('permissions', []);
                          setFieldValue('permissions', isSelected ? permissions : []);

                          // clear all permissions and set full access
                        }}
                        isSelected={fullSelected}
                      >
                        Full Access (default) - Publish to topic, Subscribe to topic, Create Topic,
                        Read Topic, Delete Topic, Destroy Topic.
                      </Checkbox>
                    </StyledFieldset>
                  </Box>
                  <Box>
                    <h2 className="mb-1 font-semibold">Custom Access</h2>
                    <StyledFieldset>
                      <Checkbox
                        {...formik.getFieldProps('custom')}
                        onChange={(isSelected) => {
                          setFullSelected(false);
                          setCustomSelected(!!isSelected);
                          // reset permissions
                          setFieldValue('permissions', []);
                        }}
                        isSelected={!!customSelected}
                      >
                        Check to grant access for each action.
                      </Checkbox>
                    </StyledFieldset>
                    <div className="mt-5 ml-5 w-full space-y-1 md:ml-10 md:w-1/2">
                      {permissions &&
                        permissions.length > 0 &&
                        permissions.map((permission: string, key: number) => (
                          <StyledFieldset key={key}>
                            <Checkbox
                              onChange={(isSelected) => {
                                setFieldValue(
                                  'permissions',
                                  fullSelected ? [] : [...values.permissions]
                                );
                                setFullSelected(false);
                                setCustomSelected(isSelected);
                                setFieldValue(
                                  'permissions',
                                  isSelected
                                    ? [...values.permissions, permission]
                                    : values.permissions.filter((p) => p !== permission)
                                );
                              }}
                              isSelected={customSelected && values.permissions.includes(permission)}
                            >
                              {permission}
                            </Checkbox>
                          </StyledFieldset>
                        ))}
                    </div>
                  </Box>
                </div>
              </fieldset>
              <div className="item-center flex  justify-center">
                <Button isLoading={isCreatingKey} className="bg-[#6DD19C] px-6 py-3 font-semibold">
                  Generate API Key
                </Button>
              </div>
            </Form>
          </div>
        </FormikProvider>
      </>
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
