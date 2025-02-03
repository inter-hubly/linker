package repository

//
// const containerEnvironments = true
//
// func TestWhatsApp(t *testing.T) {
// 	ctx := context.Background()
//
// 	var host string
// 	var close func(ctx context.Context) error
// 	var err error
//
// 	if containerEnvironments {
// 		host, close, err = testutils.ElasticSearch(ctx)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if close != nil {
// 			defer close(ctx)
// 		}
// 	} else {
// 		os.Setenv("ENVIRONMENT", "test")
// 		server.MockStartEnv(ctx, "../../")
// 		host = server.GetElasticSearch().Host
// 	}
//
// 	elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))
// 	repository := NewWhatsApp()
// 	for _, v := range []struct {
// 		testName string
// 		auxFunc  func() (string, error)
// 	}{
// 		{
// 			testName: "have to persist",
// 			auxFunc: func() (string, error) {
// 				chatToSave := testutil.GetChatToSave(testutil.NewWhatsAppMessage())
//
// 				chatToSave.Audit = append(chatToSave.Audit, entity.ChatMessageStatusTime{
// 					Status:     dto.DeliveredStatus,
// 					ReceivedAt: time.Now().Unix(),
// 				})
// 				return repository.PersistMessage(ctx, chatToSave)
// 			},
// 		},
// 		{
// 			testName: "Need change status",
// 			auxFunc: func() (string, error) {
// 				chatToSave := testutil.GetChatToSave(testutil.NewWhatsAppMessage())
//
// 				chatToSave.Audit = append(chatToSave.Audit, entity.ChatMessageStatusTime{
// 					Status:     dto.DeliveredStatus,
// 					ReceivedAt: time.Now().Unix(),
// 				})
//
// 				chatId, err := repository.PersistMessage(ctx, chatToSave)
// 				assert.Nil(t, err)
// 				assert.NotEmpty(t, chatId)
//
// 				err = repository.SetStatusMessageById(ctx, chatToSave.MessageId, dto.SentStatus, 0)
// 				assert.Nil(t, err)
//
// 				resp, err := repository.elastic.FindById(ctx, "whatsapp.ready", chatId)
// 				assert.Nil(t, err)
// 				statusResp := resp.Source["status"].([]interface{})
// 				assert.Equal(t, len(statusResp), 2)
// 				return chatId, nil
// 			},
// 		},
// 	} {
// 		t.Run(v.testName, func(t *testing.T) {
//
// 			id, err := v.auxFunc()
// 			assert.Nil(t, err)
// 			assert.NotEmpty(t, id)
// 		})
// 	}
// }
