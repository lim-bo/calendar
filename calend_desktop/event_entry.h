#ifndef EVENT_ENTRY_H
#define EVENT_ENTRY_H

#include <QWidget>
#include "eventData.h"
#include "chatwindow.h"
#include "event.h"
#include "attachments.h"
#include "allparticipants.h"
#include "replay_event.h"
namespace Ui {
class event_entry;
}

class event_entry : public QWidget
{
    Q_OBJECT

public:
    explicit event_entry(EventData data, QString viewerUID,QWidget *parent = nullptr);
    ~event_entry();
    const EventData getData() const;
private slots:
    void on_pushButton_clicked();
    void on_pushButton_3_clicked();
    void on_editButton_clicked();
    void on_pushButton_4_clicked();
    void on_checkBox_stateChanged(int arg1);
    void on_part_clicked();
    void on_pushButton_5_clicked();

private:
    Ui::event_entry *ui;
    EventData data;
    QString viewerUID;
    bool isParticipant;
    bool isAccepted;

signals:
    void deleted(event_entry*);
    void edited(event_entry*);
};

#endif // EVENT_ENTRY_H
